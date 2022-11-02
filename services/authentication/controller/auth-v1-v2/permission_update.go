package auth

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type PermissionUpdate struct {
	UserId      int        `json:"user_id"          binding:"required"`
	ServiceName string     `json:"service_name"     binding:"required"`
	MenuList    []MenuData `json:"menu_permission"  binding:"required"`
}

type MenuData struct {
	MenuId     int `json:"menu_id"     binding:"required"`
	MenuType   int `json:"menu_type"   binding:"required"`
	MenuStatus int `json:"menu_status" binding:"required"`
}

func (This PermissionUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	token := c.GetString("token")
	logger.Infof("用户权限数据修改 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.撤销操作
	if code, err := This.Cancel(token); err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.修改
	if code, err := This.UpdateMenu(token); err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ("修改用户数据成功")
}

func (This PermissionUpdate) JudgeInfo(userId int) (code int, err error) {
	// 1.该用户是否是超级管理员
	//if err = model.JudgeAuthManager(userId); err != nil {
	//	logger.Error("权限不足")
	//	return 2171, err
	//}
	// 1.检测是否是服务管理员
	if err = model.JudgeManager(userId, global.CsrServiceName); err != nil {
		logger.Error(err.Error())
		return code, err
	}
	// 2.验证用户是否存在
	if err = users.GetUser(global.DBMysql, map[string]interface{}{"id": This.UserId}); err != nil {
		logger.Info("该用户不存在")
		return 2144, errors.New("")
	}
	return
}

func (This PermissionUpdate) Cancel(token string) (code int, err error) {
	// 1.获取修改前数据
	pars := map[string]interface{}{"service_name": global.CsrServiceName, "user_id": This.UserId}
	urlRepo := fmt.Sprintf("http://%s/v2/users/data/query", utils.AuthGetServiceAddr(global.DefaultServiceName))
	headers := map[string]string{"X-Access-Token": token}
	respRepo, err := utils.PostUrl(pars, urlRepo, headers)
	if err != nil {
		logger.Error("数据查询失败")
		return 2162, errors.New("")
	}
	oldUserData := UserDataResponse{}
	err = json.Unmarshal(respRepo, &oldUserData)
	if err != nil {
		logger.Error("数据查询失败")
		return 2162, errors.New("")
	}
	if oldUserData.Code != 200 {
		logger.Error("用户与数据联系修改失败")
		return oldUserData.Code, errors.New(oldUserData.Message)
	}
	// 2.获取修改前管理数据
	oldManageMap := map[string][]int{}
	for _, oldData := range oldUserData.Data {
		if _, ok := oldManageMap[oldData.DataType]; ok {
			oldManageMap[oldData.DataType] = []int{}
		}
		if oldData.Status == 1 && oldData.Operation == "3" {
			oldManageMap[oldData.DataType] = append(oldManageMap[oldData.DataType], oldData.DataId)
		}
	}
	// 3.是否有数据被取消管理权限
	temp := false
Loop:
	for dataType, oldIds := range oldManageMap {
		for _, oldId := range oldIds {
			temp = false
			for _, menuData := range This.MenuList {
				if menuData.MenuId == oldId && dataType == fmt.Sprint(menuData.MenuType) && menuData.MenuStatus == 3 {
					temp = true
					break Loop
				}
			}
		}
	}
	// 4.撤销
	if temp {
		return This.RepoCancel(token)
	}
	return
}

type CancelRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (This PermissionUpdate) RepoCancel(token string) (code int, err error) {
	// 获取数据
	pars := map[string]int{"user_id": This.UserId}
	urlRepo := fmt.Sprintf("http://%s/v1/public/permission/cancel", utils.AuthGetServiceAddr(global.RepoServiceName))
	headers := map[string]string{"X-Access-Token": token}
	respRepo, err := utils.PostUrl(pars, urlRepo, headers)
	if err != nil {
		logger.Error("用户与数据联系修改失败")
		return 2147, errors.New("")
	}
	var cancelRes CancelRes
	err = json.Unmarshal(respRepo, &cancelRes)
	if err != nil {
		logger.Error("用户与数据联系修改失败")
		return 2147, errors.New("")
	}
	if cancelRes.Code != 200 {
		logger.Error(cancelRes.Message)
		return cancelRes.Code, errors.New(cancelRes.Message)
	}
	return
}

type UpdateMenuRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (This PermissionUpdate) UpdateMenu(token string) (code int, err error) {
	// 1.整理参数
	var dataList []map[string]interface{}
	for _, menuData := range This.MenuList {
		dataMap := map[string]interface{}{
			"data_id": menuData.MenuId, "data_type": fmt.Sprint(menuData.MenuType), "operation": fmt.Sprint(menuData.MenuStatus),
		}
		dataList = append(dataList, dataMap)
	}
	// 2.修改数据
	pars := map[string]interface{}{"service_name": global.CsrServiceName, "user_id": This.UserId, "data_list": dataList}
	urlRepo := fmt.Sprintf("http://%s/%s/users/data/update", utils.AuthGetServiceAddr(global.DefaultServiceName), global.Version)
	headers := map[string]string{"X-Access-Token": token}
	respRepo, err := utils.PostUrl(pars, urlRepo, headers)
	if err != nil {
		logger.Error("用户与数据联系修改失败")
		return 2147, errors.New("")
	}
	res := UpdateMenuRes{}
	err = json.Unmarshal(respRepo, &res)
	if err != nil {
		logger.Error("用户与数据联系修改失败")
		return 2147, errors.New("")
	}
	if res.Code != 200 {
		logger.Error(res.Message)
		return res.Code, errors.New(res.Message)
	}
	return
}

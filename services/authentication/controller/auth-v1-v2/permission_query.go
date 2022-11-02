package auth

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PermissionQuery struct {
	UserId      int    `json:"user_id"       binding:"required"`
	ServiceName string `json:"service_name"  binding:"required"`
}

type PermissionQueryRes struct {
	MenuId      int    `json:"menu_id"`
	MenuName    string `json:"menu_name"`
	MenuStatus  int    `json:"menu_status"`
	MenuType    int    `json:"menu_type"`
	MenuManager int    `json:"menu_manager"`
}

func (This PermissionQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	token := c.GetString("token")
	logger.Infof("用户权限数据查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取用户数据
	dataList, code, err := This.GetData(token)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.转换返回值
	res, code, err := This.GetRes(dataList)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 4.空数组
	if len(res) == 0 {
		res = []PermissionQueryRes{}
	}
	return ico.Succ(res)
}

func (This PermissionQuery) JudgeInfo(userId int) (code int, err error) {
	// 1.该用户是否是超级管理员
	//if err = model.JudgeAuthManager(userId); err != nil {
	//	logger.Error("权限不足")
	//	return 2171, err
	//}
	// 2.检测是否是服务管理员
	if err = model.JudgeManager(userId, global.CsrServiceName); err != nil {
		logger.Error(err.Error())
		return code, err
	}
	return
}

type UserDataResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []UserData `json:"data"`
}

type UserData struct {
	DataId     int    `json:"data_id"`
	DataName   string `json:"data_name"`
	DataType   string `json:"data_type"`
	Operation  string `json:"operation"`
	DataDesc   string `json:"data_desc"`
	Status     int    `json:"status"`
	BindStatus int    `json:"bind_status"`
}

func (This PermissionQuery) GetData(token string) (resData []UserData, code int, err error) {
	// 获取数据
	pars := map[string]interface{}{"service_name": global.CsrServiceName, "user_id": This.UserId}
	urlRepo := fmt.Sprintf("http://%s/v2/users/data/query", utils.AuthGetServiceAddr(global.DefaultServiceName))
	headers := map[string]string{"X-Access-Token": token}
	respRepo, err := utils.PostUrl(pars, urlRepo, headers)
	if err != nil {
		logger.Info("数据查询失败")
		return resData, 2162, errors.New("")
	}
	res := UserDataResponse{}
	err = json.Unmarshal(respRepo, &res)
	if err != nil {
		logger.Info("数据查询失败")
		return resData, 2162, errors.New("")
	}
	if res.Code != 200 {
		return resData, res.Code, errors.New(res.Message)
	}
	resData = res.Data
	return
}

func (This PermissionQuery) GetRes(dataList []UserData) (res []PermissionQueryRes, code int, err error) {
	// 1.排序并去重 {type: {id: {name: "1", bind: "1", status: 1}}
	dataMap := map[string]map[int]map[string]string{}
	dataOrderMap := map[string][]int{}
	for _, menu := range dataList {
		if _, ok := dataMap[menu.DataType]; !ok {
			dataMap[menu.DataType] = map[int]map[string]string{}
			dataOrderMap[menu.DataType] = []int{}
		}
		if _, ok := dataMap[menu.DataType][menu.DataId]; !ok {
			dataMap[menu.DataType][menu.DataId] = map[string]string{}
			dataOrderMap[menu.DataType] = append(dataOrderMap[menu.DataType], menu.DataId)
		}
		// 1.1 是否可选 bind
		if menu.Operation == "3" {
			if v, ok := dataMap[menu.DataType][menu.DataId]["bind"]; (ok && v == "1") || !ok {
				dataMap[menu.DataType][menu.DataId]["bind"] = fmt.Sprint(menu.BindStatus + 1)
			}
		} else {
			if _, ok := dataMap[menu.DataType][menu.DataId]["bind"]; !ok {
				dataMap[menu.DataType][menu.DataId]["bind"] = "1"
			}
		}
		// 1.2 名称 name
		dataMap[menu.DataType][menu.DataId]["name"] = menu.DataName
		// 1.3 是否绑定 status
		if menu.Status == 1 {
			dataMap[menu.DataType][menu.DataId]["status"] = menu.Operation
		} else if _, ok := dataMap[menu.DataType][menu.DataId]["status"]; !ok {
			dataMap[menu.DataType][menu.DataId]["status"] = "1"
		}
	}
	// 2.整理数据
	for dataType, dataIds := range dataOrderMap {
		menuType, _ := strconv.Atoi(dataType)
		for _, dataId := range dataIds {
			menuStatus, _ := strconv.Atoi(dataMap[dataType][dataId]["status"])
			menuManager, _ := strconv.Atoi(dataMap[dataType][dataId]["bind"])
			menuInfo := PermissionQueryRes{
				MenuId:      dataId,
				MenuName:    dataMap[dataType][dataId]["name"],
				MenuType:    menuType,
				MenuStatus:  menuStatus,
				MenuManager: menuManager,
			}
			res = append(res, menuInfo)
		}
	}
	return
}

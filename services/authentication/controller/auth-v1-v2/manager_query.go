package auth

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ManagerQuery struct {
	DataType   int   `json:"data_type" binding:"required"`
	DataIdList []int `json:"data_ids"`
}

type ManagerQueryRes struct {
	DataId      int    `json:"data_id"`
	ManagerId   int    `json:"manager_id"`
	ManagerName string `json:"manager_name"`
}

func (This ManagerQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	token := c.GetString("token")
	logger.Infof("维护用户名称查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取用户名
	userDataList, code, err := This.GetDataUser(token)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.整理返回值
	res, code, err := This.GetRes(userDataList)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 4.空数组
	if len(res) == 0 {
		res = []ManagerQueryRes{}
	}
	return ico.Succ(res)
}

func (This ManagerQuery) JudgeInfo(userId int) (code int, err error) {
	return
}

type DataUserRes struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []DataUser `json:"data"`
}

type DataUser struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	Operation string `json:"operation"`
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Account   string `json:"account"`
}

func (This ManagerQuery) GetDataUser(token string) (resData []DataUser, code int, err error) {
	// 1.整理入参
	var dataList []map[string]interface{}
	for _, dataId := range This.DataIdList {
		dataList = append(dataList, map[string]interface{}{
			"data_id": dataId, "data_type": fmt.Sprint(This.DataType), "operation": "3",
		})
	}
	pars := map[string]interface{}{"service_name": global.CsrServiceName, "data_list": dataList}
	urlRepo := fmt.Sprintf("http://%s/v2/data/user/query", utils.AuthGetServiceAddr(global.DefaultServiceName))
	headers := map[string]string{"X-Access-Token": token}
	// 2.获取数据用户
	respRepo, err := utils.PostUrl(pars, urlRepo, headers)
	if err != nil {
		logger.Info("数据查询失败")
		return resData, 2162, errors.New("")
	}
	res := DataUserRes{}
	err = json.Unmarshal(respRepo, &res)
	if err != nil {
		logger.Info("数据查询失败", string(respRepo), pars)
		return resData, 2162, errors.New("")
	}
	if res.Code != 200 {
		logger.Info("数据查询失败")
		return resData, res.Code, errors.New(res.Message)
	}
	resData = res.Data
	return
}

func (This ManagerQuery) GetRes(resDataList []DataUser) (res []ManagerQueryRes, code int, err error) {
	// 获取数据用户
	for _, resData := range resDataList {
		res = append(res, ManagerQueryRes{
			DataId: resData.DataId, ManagerId: resData.UserId, ManagerName: resData.UserName,
		})
	}
	return
}

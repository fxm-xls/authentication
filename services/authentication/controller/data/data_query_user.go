package data

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/entity"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DataUserQuery struct {
	ServiceName string         `json:"service_name" binding:"required"`
	DataList    []DataUserInfo `json:"data_list"`
}

type DataUserInfo struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	Operation string `json:"operation"`
}

type DataUserQueryRes struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	Operation string `json:"operation"`
	UserId    int    `json:"user_id"`
	Account   string `json:"account"`
	UserName  string `json:"user_name"`
}

func (This DataUserQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("数据查询_用户 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取用户
	res, code, err := This.GetDataUser()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.空数组
	if len(res) == 0 {
		res = []DataUserQueryRes{}
	}
	return ico.Succ(res)
}

func (This DataUserQuery) JudgeInfo() (code int, err error) {
	// 1.数据是否存在
	return
}

func (This DataUserQuery) GetDataUser() (res []DataUserQueryRes, code int, err error) {
	// 1.获取数据
	dataV := entity.DataView{}
	dataAllList, err := dataV.QueryListByFilter(global.DBMysql, map[string]interface{}{"service_name": This.ServiceName})
	if err != nil {
		message := fmt.Sprintf("%s", This.ServiceName)
		logger.Error("数据查询失败 ", message)
		return res, 2162, errors.New(message)
	}
	// 2.查询用户
	for _, data := range This.DataList {
		dataUserQueryRes := DataUserQueryRes{
			DataId: data.DataId, DataType: data.DataType, Operation: data.Operation,
		}
		for _, dataAll := range dataAllList {
			if data.DataId == dataAll.DataId && data.DataType == dataAll.DataType && data.Operation == dataAll.Operation {
				dataUserQueryRes.UserId = dataAll.UserId
				dataUserQueryRes.UserName = dataAll.UserName
				dataUserQueryRes.Account = dataAll.Account
			}
		}
		res = append(res, dataUserQueryRes)
	}
	return
}

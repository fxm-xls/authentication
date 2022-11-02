package data

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/services"
	"errors"
	"github.com/gin-gonic/gin"
)

type DataQuery struct {
	ServiceName string `json:"service_name" binding:"required"`
}

type DataQueryRes struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	DataName  string `json:"data_name"`
	DataDesc  string `json:"data_desc"`
	Operation string `json:"operation"`
}

func (This DataQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("数据查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.该用户是否是该服务管理员
	//if code, err := public.JudgeManager(userId, This.ServiceName); err != nil {
	//	return ico.Err(code, err.Error())
	//}
	// 2.获取单个服务
	res, code, err := This.GetData()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.空数组
	if len(res) == 0 {
		res = []DataQueryRes{}
	}
	return ico.Succ(res)
}

func (This DataQuery) GetData() (res []DataQueryRes, code int, err error) {
	// 1.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("服务不存在 ", This.ServiceName)
		return res, 2104, errors.New(This.ServiceName)
	}
	// 2.获取数据
	entityT := entity.EntityTable{}
	entieyList, err := entityT.QueryListByFilter(global.DBMysql, map[string]interface{}{"service_id": serviceId})
	if err != nil {
		logger.Error("数据查询失败 ", This.ServiceName)
		return res, 2162, errors.New(This.ServiceName)
	}
	for _, entityInfo := range entieyList {
		res = append(res, DataQueryRes{
			DataId: entityInfo.DataId, DataType: entityInfo.DataType, Operation: entityInfo.Operation,
			DataName: entityInfo.DataName, DataDesc: entityInfo.DataDesc,
		})
	}
	return
}

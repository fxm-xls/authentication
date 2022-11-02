package services

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/services"
	"errors"
	"github.com/gin-gonic/gin"
)

type ServiceQuery struct{}

type ServiceQueryRes struct {
	ServiceId   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	ServiceDesc string `json:"service_desc"`
	IndexUrl    string `json:"index_url"`
}

func (This ServiceQuery) DoHandle(c *gin.Context) *ico.Result {
	logger.Infof("服务查询 user: id %d, name %s", c.GetInt("user_id"), c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.查询
	res, code, err := This.GetService()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.空数组
	if len(res) == 0 {
		res = []ServiceQueryRes{}
	}
	return ico.Succ(res)
}

func (This ServiceQuery) JudgeInfo() (code int, err error) {
	return
}

func (This ServiceQuery) GetService() (res []ServiceQueryRes, code int, err error) {
	// 1.获取服务数据
	serviceT := services.ServiceTable{}
	filterMap := map[string]interface{}{}
	serviceList, err := serviceT.QueryListByFilter(global.DBMysql, filterMap)
	if err != nil {
		return res, 2101, errors.New("")
	}
	for _, service := range serviceList {
		res = append(res, ServiceQueryRes{
			ServiceId: service.ServiceId, ServiceName: service.ServiceName,
			ServiceDesc: service.ServiceDesc, IndexUrl: service.IndexUrl,
		})
	}
	return
}

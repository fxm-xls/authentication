package interfaces

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/interfaces"
	"bigrule/services/authentication/model/services"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type InterfaceQuery struct {
	ServiceName string `json:"service_name" binding:"required"`
}

func (This InterfaceQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("接口查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.查询服务下接口
	res, code, err := This.GetInterfaces()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.空数组
	if len(res) == 0 {
		res = []interfaces.InterfaceTable{}
	}
	return ico.Succ(res)
}

func (This InterfaceQuery) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if code, err = public.JudgeManager(userId, This.ServiceName); err != nil {
		logger.Info("权限不足")
		return 2171, errors.New("")
	}
	// 2.验证service-name是否规范
	if This.ServiceName != "-1" {
		if !strings.HasSuffix(This.ServiceName, "-service") || This.ServiceName == "-service" {
			logger.Info("服务名称格式异常，服务名称需以-service结尾")
			return 2101, errors.New("服务名称需以-service结尾")
		}
	}
	return
}

func (This InterfaceQuery) GetInterfaces() (res []interfaces.InterfaceTable, code int, err error) {
	// 0.服务名筛选
	whereMap := map[string]interface{}{}
	if This.ServiceName != "-1" {
		serviceId, err := services.GetServiceIdByName(This.ServiceName)
		if err != nil {
			logger.Info("该服务不存在")
			return res, 2104, errors.New(This.ServiceName)
		}
		whereMap["service_id"] = serviceId
	}
	interfaceT := interfaces.InterfaceTable{}
	// 1.查询
	interfaceList, err := interfaceT.QueryListByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Info("接口查询失败")
		return res, 2152, errors.New("")
	}
	res = interfaceList
	return
}

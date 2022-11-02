package services

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

type ServiceUpdate struct {
	ServiceId   int    `json:"service_id"   binding:"required"`
	ServiceName string `json:"service_name"`
	IndexUrl    string `json:"index_url"`
	ServiceDesc string `json:"service_desc"`
}

func (This ServiceUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	logger.Infof("服务修改 user: id %d, name %s", c.GetInt("user_id"), c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取旧服务名
	serviceT := services.ServiceTable{}
	service, err := serviceT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.ServiceId})
	if err != nil {
		return ico.Err(2101, "")
	}
	// 3.修改基本信息
	if code, err := This.UpdateInfo(service.ServiceId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 4.服务名不需修改
	if service.ServiceName == This.ServiceName || This.ServiceName == "" {
		return ico.Succ("修改服务成功")
	}
	// 5.修改内置管理员角色-用户_名称
	tx := global.DBMysql.Begin()
	if code, err := This.UpdateDefault(tx, service.ServiceName); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("修改服务成功")
}

func (This ServiceUpdate) JudgeInfo() (code int, err error) {
	if This.ServiceName != "" {
		// 需以 -service结尾
		if !strings.HasSuffix(This.ServiceName, "-service") || This.ServiceName == "-service" {
			logger.Info("服务名称格式异常")
			return 2101, errors.New("")
		}
	}
	return
}

func (This ServiceUpdate) UpdateInfo(serviceId int) (code int, err error) {
	serviceT := services.ServiceTable{
		ServiceId: serviceId, ServiceName: This.ServiceName, ServiceDesc: This.ServiceDesc, IndexUrl: This.IndexUrl,
	}
	if err := serviceT.UpdateByStruct(global.DBMysql); err != nil {
		logger.Info("服务修改失败")
		return 2103, errors.New("")
	}
	return
}

func (This ServiceUpdate) UpdateDefault(tx *gorm.DB, serviceName string) (code int, err error) {
	// 1.修改内置角色
	roleT := roles.RoleTable{}
	whereMap := map[string]interface{}{"service_id": This.ServiceId, "manager": global.DefaultManagerInt}
	roleList, err := roleT.QueryListByFilter(tx, whereMap)
	if err != nil || len(roleList) < 1 {
		logger.Info("角色查询失败")
		return 2112, errors.New("")
	}
	upInfo := map[string]interface{}{
		"name": fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix),
		"desc": fmt.Sprintf("%s 内置管理员", This.ServiceName),
	}
	if err := roleList[0].Update(tx, upInfo); err != nil {
		logger.Info("角色修改失败")
		return 2111, errors.New("")
	}
	// 2.修改内置账户
	userT := users.UserTable{}
	whereMap = map[string]interface{}{"account": fmt.Sprintf("%s_%s", serviceName, global.DefaultSuffix)}
	userList, err := userT.QueryListByFilter(tx, whereMap)
	if err != nil || len(userList) < 1 {
		logger.Info("用户查询失败")
		return 2135, errors.New("")
	}
	password, _ := utils.HashPassword(fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix))
	upInfo = map[string]interface{}{
		"account":   fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix),
		"user_name": fmt.Sprintf("%s_内置管理员账户", This.ServiceName),
		"desc":      fmt.Sprintf("%s_内置管理员账户", This.ServiceName),
		"password":  password,
	}
	if err := userList[0].Update(tx, upInfo); err != nil {
		logger.Info("用户修改失败")
		return 2131, errors.New("")
	}
	return
}

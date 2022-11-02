package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type RoleAdd struct {
	RoleName    string `json:"role_name"     binding:"required"`
	RoleDesc    string `json:"role_desc"     binding:"required"`
	ServiceName string `json:"service_name"  binding:"required"`
	UserNum     int    `json:"user_num"`
	RoleNum     int    `json:"role_num"`
}

type RoleAddRes struct {
	RoleId int `json:"role_id"`
}

func (This RoleAdd) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色新增-基本信息 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("该服务不存在 ", This.ServiceName)
		return ico.Err(2104, This.ServiceName)
	}
	// 3.增加角色
	roleT := roles.RoleTable{
		RoleName:   This.RoleName,
		RoleDesc:   This.RoleDesc,
		ServiceId:  serviceId,
		CreateTime: time.Now().Unix(),
		UserNum:    This.UserNum,
		RoleNum:    This.RoleNum,
	}
	if This.UserNum <= 0 {
		roleT.UserNum = global.DefaultLimitNum
	}
	if This.RoleNum <= 0 {
		roleT.RoleNum = global.DefaultLimitNum
	}
	if err = roleT.Insert(global.DBMysql); err != nil {
		return ico.Err(2111, "")
	}
	res := RoleAddRes{RoleId: roleT.RoleId}
	return ico.Succ(res)
}

func (This RoleAdd) JudgeInfo(userId int) (code int, err error) {
	// 1.该用户是否是超级管理员
	if err = model.JudgeAuthManager(userId); err != nil {
		logger.Error("权限不足")
		return 2171, err
	}
	// 2.角色名称不能为空字符
	if strings.TrimSpace(This.RoleName) == "" || strings.TrimSpace(This.RoleDesc) == "" {
		logger.Error("角色名称不能为空字符")
		return 2129, errors.New("名称或描述不能为空字符")
	}
	// 3.该角色是否已存在
	roleT := roles.RoleTable{}
	if _, err = roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"name": This.RoleName}); err == nil {
		logger.Error("角色已存在")
		return 2128, errors.New("")
	}
	err = nil
	return
}

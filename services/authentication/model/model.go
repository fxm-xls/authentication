package model

import (
	"bigrule/common/global"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"errors"
)

// JudgeAuthManager 判断用户是否为超级管理员
func JudgeAuthManager(userId int) (err error) {
	// 1.获取超级管理员角色
	roleId := global.AuthManagerId
	// 2.该用户是否配置角色
	whereMap := map[string]interface{}{"user_id": userId, "role_id": roleId}
	userRoleT := users.UserRoleTable{}
	if _, err = userRoleT.QueryByFilter(global.DBMysql, whereMap); err == nil {
		return
	}
	// 3.该用户是否包含角色
	userRoleV := users.UserRoleView{}
	if _, err = userRoleV.QueryByFilter(global.DBMysql, whereMap); err != nil {
		logger.Error("权限不足")
		return
	}
	return
}

// JudgeManager 判断用户是否为该服务管理员
func JudgeManager(userId int, serviceName string) (err error) {
	// 1.获取该服务id
	serviceId, err := services.GetServiceIdByName(serviceName)
	if err != nil {
		logger.Error(serviceName, "该服务不存在")
		return errors.New("该服务不存在")
	}
	// 2.获取该服务管理员角色
	roleT := roles.RoleTable{}
	whereMap := map[string]interface{}{"service_id": serviceId, "manager": global.DefaultManagerInt}
	role, err := roleT.QueryByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Error(serviceName, "该服务不存在")
		return errors.New("该服务不存在")
	}
	// 3.该用户是否配置角色
	whereMap = map[string]interface{}{"user_id": userId, "role_id": role.RoleId}
	userRoleT := users.UserRoleTable{}
	if _, err = userRoleT.QueryByFilter(global.DBMysql, whereMap); err == nil {
		return
	}
	// 4.该用户是否包含角色
	userRoleV := users.UserRoleView{}
	if _, err = userRoleV.QueryByFilter(global.DBMysql, whereMap); err != nil {
		logger.Error("权限不足")
		return
	}
	return
}

// JudgeDepartmentManager 判断用户是否为部门管理员
func JudgeDepartmentManager(userId int) (err error) {
	departmentT := departments.DepartmentTab{}
	whereMap := map[string]interface{}{"charger_id": userId}
	_, err = departmentT.QueryByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Error("权限不足")
		return errors.New("权限不足")
	}
	return
}

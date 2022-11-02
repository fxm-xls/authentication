package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/users"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleDelete struct {
	RoleId      int    `json:"role_id"       binding:"required"`
	ServiceName string `json:"service_name"  binding:"required"`
}

func (This RoleDelete) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色删除 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.删除
	var tx = global.DBMysql.Begin()
	if code, err := This.DeleteRole(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("删除成功")
}

func (This RoleDelete) JudgeInfo(userId int) (code int, err error) {
	// 1.该用户是否是超级管理员
	if err = model.JudgeAuthManager(userId); err != nil {
		logger.Error("权限不足")
		return 2171, err
	}
	// 2.检测角色是否存在
	roleT := roles.RoleTable{}
	role, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId})
	if err != nil {
		logger.Error("角色不存在")
		return 2119, errors.New("")
	}
	// 3.检测角色是否是内置
	if role.Manager == 1 {
		logger.Error("管理员角色不可删除")
		return 2126, errors.New("")
	}
	// 4.检测角色是否被用户使用
	userRoleT := users.UserRoleTable{}
	if _, err = userRoleT.QueryByFilter(global.DBMysql, map[string]interface{}{"role_id": This.RoleId}); err == nil {
		logger.Error("角色仍有用户使用")
		return 2125, errors.New("")
	}
	// 5.检测角色是否被部门使用
	deptRoleT := departments.DeptRoleTab{}
	if _, err = deptRoleT.QueryByFilter(global.DBMysql, map[string]interface{}{"role_id": This.RoleId}); err == nil {
		logger.Error("角色仍有部门使用")
		return 2127, errors.New("")
	}
	return 200, nil
}

func (This RoleDelete) DeleteRole(tx *gorm.DB) (code int, err error) {
	// 1.删除角色信息
	roleT := roles.RoleTable{RoleId: This.RoleId}
	if err = roleT.Delete(tx); err != nil {
		logger.Error("角色信息删除失败")
		return 2116, errors.New("")
	}
	// 2.删除对应接口联系
	roleInterfaceT := roles.RoleInterfaceTable{}
	if err = roleInterfaceT.DeleteByRoleIds(tx, []int{This.RoleId}); err != nil {
		logger.Error("角色与接口联系删除失败")
		return 2113, errors.New("")
	}
	// 3.删除对应角色联系
	roleSubT := roles.RoleSubTable{}
	if err = roleSubT.DeleteBySubRoleIds(tx, []int{This.RoleId}); err != nil {
		logger.Error("角色与大角色联系删除失败")
		return 2136, errors.New("")
	}
	if err = roleSubT.DeleteByRoleIds(tx, []int{This.RoleId}); err != nil {
		logger.Error("角色与大角色联系删除失败")
		return 2136, errors.New("")
	}
	// 4.删除对应casbin
	e := mycasbin.Casbin()
	if _, err = e.RemoveFilteredPolicy(0, fmt.Sprint(This.RoleId)); err != nil {
		logger.Error("角色与casbin联系删除失败")
		return 2115, errors.New("")
	}
	return
}

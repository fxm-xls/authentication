package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/roles"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleSubUpdate struct {
	RoleId     int   `json:"role_id"       binding:"required"`
	SubRoleIds []int `json:"sub_role_ids"  binding:"required"`
}

func (This RoleSubUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色修改-大角色分配小角色 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.为角色分配小角色
	tx := global.DBMysql.Begin()
	if code, err := This.UpdateRole(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("修改成功")
}

func (This RoleSubUpdate) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if err = model.JudgeAuthManager(userId); err != nil {
		logger.Error("权限不足")
		return 2171, errors.New("")
	}
	// 2.查询大角色是否存在
	roleT := roles.RoleTable{}
	if _, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId, "service_id": global.DefaultServiceId}); err != nil {
		logger.Error("角色不存在")
		return 2119, errors.New("")
	}
	// 3.查询小角色是否存在
	roleList, err := roleT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error("角色不存在")
		return 2119, errors.New("")
	}
	for _, subRoleId := range This.SubRoleIds {
		temp := false
		for _, role := range roleList {
			if subRoleId == role.RoleId && role.ServiceId != global.DefaultServiceId {
				temp = true
				break
			}
		}
		if !temp {
			logger.Error("角色不存在")
			return 2119, errors.New("小角色")
		}
	}
	return
}

func (This RoleSubUpdate) UpdateRole(tx *gorm.DB) (code int, err error) {
	roleSubT := roles.RoleSubTable{}
	// 1.先删除
	if err = roleSubT.DeleteByRoleIds(tx, []int{This.RoleId}); err != nil {
		logger.Error("角色与大角色联系删除失败")
		return 2114, errors.New("")
	}
	// 2.再增加
	var roleSubList []roles.RoleSubTable
	for _, subRoleId := range This.SubRoleIds {
		roleSubList = append(roleSubList, roles.RoleSubTable{RoleId: This.RoleId, SubRoleId: subRoleId})
	}
	if len(roleSubList) > 0 {
		if err = roleSubT.InsertMany(tx, roleSubList); err != nil {
			logger.Error("角色与大角色联系修改失败")
			return 2122, errors.New("")
		}
	}
	return
}

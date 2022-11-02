package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/roles"
	"github.com/gin-gonic/gin"
)

type RoleUpdate struct {
	RoleId      int    `json:"role_id"       binding:"required"`
	RoleName    string `json:"role_name"`
	RoleDesc    string `json:"role_desc"`
	UserNum     int    `json:"user_num"`
	RoleNum     int    `json:"role_num"`
}

func (This RoleUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色修改-基本信息 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.该用户是否是超级管理员
	if err := model.JudgeAuthManager(userId); err != nil {
		return ico.Err(2171, err.Error())
	}
	roleT := roles.RoleTable{
		RoleId: This.RoleId, RoleName: This.RoleName, RoleDesc: This.RoleDesc,
		UserNum: This.UserNum, RoleNum: This.RoleNum,
	}
	// 2.查询角色是否存在
	if This.RoleId != 0 {
		if _, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId}); err != nil {
			logger.Info("角色不存在")
			return ico.Err(2119, "")
		}
	}
	// 3.修改角色
	if err := roleT.UpdateByStruct(global.DBMysql); err != nil {
		logger.Info("角色基本信息修改失败")
		return ico.Err(2120, "")
	}
	return ico.Succ("修改成功")
}

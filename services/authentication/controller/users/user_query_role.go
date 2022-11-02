package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/users"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserRoleQuery struct {
	UserId int `json:"user_id"       binding:"required"`
}

type UserRoleQueryRes struct {
	Account  string     `json:"account"`
	UserName string     `json:"user_name"`
	RoleList []RoleInfo `json:"role_list"`
}

type RoleInfo struct {
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	Status   int    `json:"status"`
}

func (This UserRoleQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("用户查询_角色 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取用户
	res, code, err := This.GetUsers()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ(res)
}

func (This UserRoleQuery) JudgeInfo(userId int) (code int, err error) {
	// 1.检测账户是否存在
	userT := users.UserTable{}
	if _, err = userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.UserId}); err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	// 2.管理员检测
	if err = model.JudgeAuthManager(userId); err != nil {
		logger.Error("权限不足")
		return 2171, errors.New("")
	}
	return
}

func (This UserRoleQuery) GetUsers() (res UserRoleQueryRes, code int, err error) {
	// 1.获取Auth所有角色
	roleT := roles.RoleTable{}
	roleList, err := roleT.QueryListByFilter(global.DBMysql, map[string]interface{}{"service_id": global.DefaultServiceId})
	if err != nil {
		logger.Error("角色查询失败")
		return res, 2112, errors.New("")
	}
	// 2.获取该用户角色
	userRoleT := users.UserRoleTable{}
	userRoleList, err := userRoleT.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": This.UserId})
	if err != nil {
		logger.Error("角色查询失败")
		return res, 2112, errors.New("")
	}
	// 3.配置
	for _, role := range roleList {
		roleInfo := RoleInfo{RoleId: role.RoleId, RoleName: role.RoleName}
		for _, userRole := range userRoleList {
			if userRole.RoleId == role.RoleId {
				roleInfo.Status = 1
				break
			}
		}
		res.RoleList = append(res.RoleList, roleInfo)
	}
	return
}

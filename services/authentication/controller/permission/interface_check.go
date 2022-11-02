package permission

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/middleware"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/middleware/jwt"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type InterfaceCheck struct {
	Path   string `json:"path"   binding:"required"`
	Method string `json:"method" binding:"required"`
}

func (This InterfaceCheck) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "")
	}
	// 0.校验jwt path
	if code, err := jwt.ManageAuthToken(This.Path, c); err != nil {
		return ico.ErrJwt(code, err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("用户权限验证 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.获取roleIds
	roleIds, code, err := This.GetRole(userId)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.过滤普通接口
	if code, err = This.JudgeInterface(userId); err == nil {
		return ico.Succ("True")
	}
	// 3.casbin校验
	if err := mycasbin.AuthCheckRole(roleIds, This.Path, This.Method); err != nil {
		return ico.Err(2177, "")
	}
	return ico.Succ("True")
}

func (This InterfaceCheck) GetRole(userId int) (roleIds []int, code int, err error) {
	userRoleT := users.UserRoleTable{}
	userRoleList, err := userRoleT.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": userId})
	if err != nil {
		logger.Info("角色查询失败")
		return roleIds, 2112, errors.New(fmt.Sprint("userId:", userId))
	}
	for _, userRole := range userRoleList {
		roleIds = append(roleIds, userRole.RoleId)
	}
	userRoleV := users.UserRoleView{}
	userRoleVList, err := userRoleV.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": userId})
	if err != nil {
		logger.Info("角色查询失败")
		return roleIds, 2112, errors.New(fmt.Sprint("userId:", userId))
	}
	for _, userRole := range userRoleVList {
		roleIds = append(roleIds, userRole.RoleId)
	}
	return
}

func (This InterfaceCheck) JudgeInterface(userId int) (code int, err error) {
	// 1.去除普通接口
	if utils.IsContains(middleware.CasbinNoVerify, This.Path) {
		logger.Info(This.Path, " 不校验")
		return
	}
	// 2.管理员默认拥有该服务所有权限
	userV := users.UserView{}
	userVList, err := userV.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": userId, "manager": global.DefaultManagerInt})
	if err != nil {
		logger.Info("角色查询失败")
		return 2112, errors.New("")
	}
	for _, user := range userVList {
		if strings.HasPrefix(This.Path, fmt.Sprintf("/%s", user.ServiceName)) {
			logger.Info(This.Path, " 不校验")
			return
		}
		if utils.IsContains(middleware.ManagerCasbin, This.Path) {
			logger.Info(This.Path, " 不校验")
			return
		}
	}
	userRoleV := users.UserRoleView{}
	userRoleVList, err := userRoleV.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": userId, "manager": global.DefaultManagerInt})
	if err != nil {
		logger.Info("角色查询失败")
		return 2112, errors.New("")
	}
	logger.Info("userVList:", userVList)
	for _, user := range userRoleVList {
		if strings.HasPrefix(This.Path, fmt.Sprintf("/%s", user.ServiceName)) {
			logger.Info(This.Path, " 不校验")
			return
		}
		if utils.IsContains(middleware.ManagerCasbin, This.Path) {
			logger.Info(This.Path, " 不校验")
			return
		}
		// 3.csr
		if user.ServiceName == "flowcsr-service" && user.Manager == 1 {
			if strings.HasPrefix(This.Path, fmt.Sprintf("/repo-bff-service")) || strings.HasPrefix(This.Path,
				fmt.Sprintf("/repo-service")) || strings.HasPrefix(This.Path, fmt.Sprintf("/compile-service")) {
				logger.Info(This.Path, " 不校验")
				return
			}
		}
	}
	return 2177, errors.New("")
}

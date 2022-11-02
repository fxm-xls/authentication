package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type RoleImport struct {
	RoleList    []Role `json:"role_list"     binding:"required"`
	ServiceName string `json:"service_name"  binding:"required"`
	ServiceId   int
}

type Role struct {
	RoleName      string      `json:"role_name"     binding:"required"`
	InterfaceList []Interface `json:"interface_list"`
}

type Interface struct {
	InterfaceName string `json:"interface_name"     binding:"required"`
	InterfaceDesc string `json:"interface_desc"`
	Path          string `json:"path"               binding:"required"`
	Method        string `json:"method"             binding:"required"`
}

type RoleImportRes struct {
	RoleId int `json:"role_id"`
}

func (This RoleImport) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色导入 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.增加角色
	tx := global.DBMysql.Begin()
	if code, err := This.ImportRole(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	return ico.Succ("")
}

func (This *RoleImport) JudgeInfo(userId int) (code int, err error) {
	// 1.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("该服务不存在 ", This.ServiceName)
		return 2104, err
	}
	This.ServiceId = serviceId
	return
}

func (This RoleImport) ImportRole(tx *gorm.DB) (code int, err error) {
	roleT := roles.RoleTable{}
	e := mycasbin.Casbin()
	// 1.查找角色
	roleList, err := roleT.QueryListByFilter(tx, map[string]interface{}{"service_id": This.ServiceId})
	if err != nil {
		return 2112, errors.New("")
	}
	// 2.删除角色
	for _, role := range roleList {
		// 2.1 去除管理员
		if role.Manager == 1 {
			continue
		}
		// 2.2 删除角色与casbin联系
		_, err = e.RemoveFilteredPolicy(0, fmt.Sprint(role.RoleId))
		if err != nil {
			logger.Info("角色与casbin联系删除失败")
			return 2115, errors.New("")
		}
		// 2.3 删除角色
		roleDeleteT := roles.RoleTable{RoleId: role.RoleId}
		if err = roleDeleteT.Delete(tx); err != nil {
			logger.Error("角色信息删除失败")
			return 2116, errors.New("")
		}
	}
	// 3.增加角色
	for _, role := range This.RoleList {
		// 3.1 增加角色
		roleAddT := roles.RoleTable{
			RoleName:   role.RoleName,
			ServiceId:  This.ServiceId,
			CreateTime: time.Now().Unix(),
			UserNum:    global.DefaultLimitNum,
			RoleNum:    global.DefaultLimitNum,
		}
		if err = roleAddT.Insert(global.DBMysql); err != nil {
			logger.Error("角色新增-基本信息失败")
			return 2111, errors.New("")
		}
		// 3.2 增加角色与casbin联系
		var rules [][]string
		for _, interfaceInfo := range role.InterfaceList {
			rules = append(rules, append([]string{fmt.Sprint(roleAddT.RoleId), interfaceInfo.Path, interfaceInfo.Method}))
		}
		_, err = e.AddNamedPolicies("p", rules)
		if err != nil {
			logger.Info("角色与casbin联系修改失败")
			return 2123, errors.New("")
		}
	}
	return
}

package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/model/interfaces"
	"bigrule/services/authentication/model/roles"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleInterfaceUpdate struct {
	RoleId       int    `json:"role_id"       binding:"required"`
	ServiceName  string `json:"service_name"  binding:"required"`
	InterfaceIds []int  `json:"interface_ids" binding:"required"`
}

func (This RoleInterfaceUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色修改-小角色分配接口 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	roleInterfaceList, newInterfaceList, code, err := This.JudgeInfo(userId)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.修改角色
	tx := global.DBMysql.Begin()
	if code, err = This.UpdateRole(tx, roleInterfaceList, newInterfaceList); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("修改成功")
}

func (This RoleInterfaceUpdate) JudgeInfo(userId int) (roleInterfaceList []roles.RoleInterfaceTable, newInterfaceList []interfaces.InterfaceTable, code int, err error) {
	// 1.该用户是否是该服务管理员
	//if code, err = public.JudgeManager(userId, This.ServiceName); err != nil {
	//	return roleInterfaceList, newInterfaceList, code, err
	//}
	// 2.查询信息是否真实
	// 2.1查询角色是否存在
	roleT := roles.RoleTable{}
	if _, err = roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId}); err != nil {
		logger.Info("角色不存在")
		return roleInterfaceList, newInterfaceList, 2119, errors.New("")
	}
	// 2.2查询接口是否存在
	interfaceT := interfaces.InterfaceTable{}
	interfaceList, err := interfaceT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Info("接口查询失败")
		return roleInterfaceList, newInterfaceList, 2152, errors.New("")
	}
	for _, interfaceId := range This.InterfaceIds {
		temp := false
		for _, interfaceInfo := range interfaceList {
			if interfaceId == interfaceInfo.InterfaceId {
				temp = true
				newInterfaceList = append(newInterfaceList, interfaceInfo)
				break
			}
		}
		if !temp {
			logger.Info("接口不存在")
			return roleInterfaceList, newInterfaceList, 2155, errors.New("")
		}
		roleInterfaceList = append(roleInterfaceList, roles.RoleInterfaceTable{
			RoleId: This.RoleId, InterfaceId: interfaceId,
		})
	}
	return
}

func (This RoleInterfaceUpdate) UpdateRole(tx *gorm.DB, roleInterfaceList []roles.RoleInterfaceTable, newInterfaceList []interfaces.InterfaceTable) (code int, err error) {
	// 1.为角色分配接口
	roleInterfaceT := roles.RoleInterfaceTable{}
	if err := roleInterfaceT.DeleteByRoleIds(tx, []int{This.RoleId}); err != nil {
		return 2113, errors.New("角色与接口联系删除失败")
	}
	if len(roleInterfaceList) > 0 {
		if err = roleInterfaceT.InsertMany(tx, roleInterfaceList); err != nil {
			logger.Info("角色与接口联系修改失败")
			return 2121, errors.New("")
		}
	}
	// 2.为角色分配casbin联系
	e := mycasbin.Casbin()
	_, err = e.RemoveFilteredPolicy(0, fmt.Sprint(This.RoleId))
	if err != nil {
		logger.Info("角色与casbin联系删除失败")
		return 2115, errors.New("")
	}
	var rules [][]string
	for _, interfaceInfo := range newInterfaceList {
		rules = append(rules, append([]string{fmt.Sprint(This.RoleId), interfaceInfo.Path, interfaceInfo.Method}))
	}
	_, err = e.AddNamedPolicies("p", rules)
	if err != nil {
		logger.Info("角色与casbin联系修改失败")
		return 2123, errors.New("")
	}
	return
}

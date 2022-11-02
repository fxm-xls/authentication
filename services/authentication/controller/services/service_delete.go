package services

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/interfaces"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceDelete struct {
	ServiceId   int `json:"service_id" binding:"required"`
	ServiceName string
}

func (This ServiceDelete) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	logger.Infof("服务删除 user: id %d, name %s", c.GetInt("user_id"), c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取服务名
	serviceT := services.ServiceTable{}
	service, err := serviceT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.ServiceId})
	if err != nil {
		logger.Info("服务不存在")
		return ico.Err(2104, "")
	}
	This.ServiceName = service.ServiceName
	// 3.删除
	tx := global.DBMysql.Begin()
	if code, err := This.DeleteDefault(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("删除服务成功")
}

func (This ServiceDelete) JudgeInfo() (code int, err error) {
	// 1.权限服务不可删除
	if This.ServiceId == global.DefaultServiceId {
		logger.Info("权限服务不可删除")
		return 2105, errors.New("")
	}
	return
}

func (This ServiceDelete) DeleteDefault(tx *gorm.DB) (code int, err error) {
	serviceT := services.ServiceTable{ServiceId: This.ServiceId}
	// 1.删除服务
	if err = serviceT.Delete(tx); err != nil {
		logger.Info("服务删除失败")
		return 2107, errors.New("")
	}
	// 2.删除数据相关信息
	if code, err = This.DeleteData(tx); err != nil {
		logger.Info("服务删除失败")
		return code, err
	}
	// 3.删除接口相关信息
	if code, err = This.DeleteInterface(tx); err != nil {
		logger.Info("服务删除失败")
		return code, err
	}
	// 4.删除角色相关信息
	if code, err = This.DeleteRole(tx); err != nil {
		logger.Info("服务删除失败")
		return code, err
	}
	// 5.删除用户相关信息
	//if code, err = This.DeleteUser(tx); err != nil {
	//	logger.Info("服务删除失败")
	//	return code, err
	//}
	return
}

func (This ServiceDelete) DeleteData(tx *gorm.DB) (code int, err error) {
	entityT := entity.EntityTable{}
	// 0.查询所有数据
	var entityIds []int
	entityList, err := entityT.QueryListByFilter(tx, map[string]interface{}{"service_id": This.ServiceId})
	if err != nil {
		logger.Info("数据查询失败")
		return 2162, errors.New("")
	}
	for _, entityInfo := range entityList {
		entityIds = append(entityIds, entityInfo.EntityId)
	}
	// 1.删除用户与该数据联系
	userEntityT := users.UserEntityTable{}
	if err = userEntityT.DeleteByEntityIds(tx, entityIds); err != nil {
		logger.Info("用户与数据联系删除失败")
		return 2137, errors.New("")
	}
	// 2.删除数据
	if err = entityT.DeleteByServiceId(tx, This.ServiceId); err != nil {
		logger.Info("数据删除失败")
		return 2163, errors.New("")
	}
	return
}

func (This ServiceDelete) DeleteUser(tx *gorm.DB) (code int, err error) {
	// 1.删除默认用户
	userT := users.UserTable{}
	whereMap := map[string]interface{}{"account": fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix)}
	if err = userT.DeleteByFilter(tx, whereMap); err != nil {
		logger.Info("接口删除失败")
		return 2153, errors.New("")
	}
	return
}

func (This ServiceDelete) DeleteRole(tx *gorm.DB) (code int, err error) {
	// 0.查询所有角色
	var roleIds []int
	roleT := roles.RoleTable{ServiceId: This.ServiceId}
	roleList, err := roleT.QueryListByFilter(tx, map[string]interface{}{"service_id": This.ServiceId})
	if err != nil || len(roleList) < 1 {
		logger.Info("角色查询失败")
		return 2112, errors.New("")
	}
	for _, role := range roleList {
		roleIds = append(roleIds, role.RoleId)
	}
	// 1.删除角色与接口联系
	roleInterfaceT := roles.RoleInterfaceTable{}
	if err = roleInterfaceT.DeleteByRoleIds(tx, roleIds); err != nil {
		logger.Info("角色与接口联系删除失败")
		return 2113, errors.New("")
	}
	// 2.删除角色与大角色联系
	roleSubT := roles.RoleSubTable{}
	if err = roleSubT.DeleteBySubRoleIds(tx, roleIds); err != nil {
		logger.Info("角色与大角色联系删除失败")
		return 2114, errors.New("")
	}
	// 3.删除账户与角色联系
	userRoleT := users.UserRoleTable{}
	if err = userRoleT.DeleteByRoleIds(tx, roleIds); err != nil {
		logger.Info("用户与角色联系删除失败")
		return 2136, errors.New("")
	}
	// 4.删除casbin角色
	e := mycasbin.Casbin()
	for _, roleId := range roleIds {
		_, err = e.RemoveFilteredPolicy(0, fmt.Sprint(roleId))
		if err != nil {
			logger.Info("角色与casbin联系删除失败")
			return 2115, errors.New("")
		}
	}
	// 5.删除角色
	if err = roleT.DeleteByRoleIds(tx, roleIds); err != nil {
		logger.Info("角色删除失败")
		return 2116, errors.New("")
	}
	return
}

func (This ServiceDelete) DeleteInterface(tx *gorm.DB) (code int, err error) {
	// 1.删除所有接口
	interfaceT := interfaces.InterfaceTable{}
	if err = interfaceT.DeleteByServiceId(tx, This.ServiceId); err != nil {
		logger.Info("接口删除失败")
		return 2153, errors.New("")
	}
	return
}

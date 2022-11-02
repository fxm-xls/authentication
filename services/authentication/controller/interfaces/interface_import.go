package interfaces

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	mycasbin "bigrule/services/authentication/middleware/casbin"
	"bigrule/services/authentication/model/interfaces"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InterfaceImport struct {
	InterfaceList  []InterfaceInfo `json:"interface_list"  binding:"required"`
	ServiceName    string          `json:"service_name"    binding:"required"`
	IsDisassociate int             `json:"is_disassociate"`
	ServiceId      int
}

type InterfaceInfo struct {
	InterfaceName string `json:"interface_name" binding:"required"`
	InterfaceDesc string `json:"interface_desc"`
	Path          string `json:"path"           binding:"required"`
	Method        string `json:"method"         binding:"required"`
}

func (This InterfaceImport) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("接口导入 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 1.1查询服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		message := fmt.Sprintf("%s", This.ServiceName)
		logger.Error("服务不存在 ", message)
		return ico.Err(2104, message)
	}
	This.ServiceId = serviceId
	// 2.查询出初始接口列表
	interfaceList, code, err := This.GetInterfaceByName()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.与本次对比
	if code, err = This.DiffInterface(interfaceList); err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ("导入成功")
}

func (This InterfaceImport) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if code, err = public.JudgeManager(userId, This.ServiceName); err != nil {
		logger.Error("权限不足")
		return 2171, errors.New("")
	}
	// 2.验证path、method、name是否存在
	for _, interfaceInfo := range This.InterfaceList {
		if interfaceInfo.InterfaceName == "" || interfaceInfo.Path == "" || interfaceInfo.Method == "" {
			logger.Error("接口信息异常")
			return 2156, errors.New("接口名称、路径、方式不能为空")
		}
	}
	// 3.验证IsDisassociate
	if This.IsDisassociate != 0 && This.IsDisassociate != 1 {
		logger.Error("参数异常")
		return 2099, errors.New("is_disassociate只能为0或1")
	}
	return
}

func (This InterfaceImport) GetInterfaceByName() (res []interfaces.InterfaceTable, code int, err error) {
	interfaceT := interfaces.InterfaceTable{}
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("该服务不存在 ", This.ServiceName)
		return res, 2104, errors.New(This.ServiceName)
	}
	This.ServiceId = serviceId
	whereMap := map[string]interface{}{"service_id": serviceId}
	res, err = interfaceT.QueryListByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Error("接口查询失败")
		return res, 2152, errors.New("")
	}
	return
}

func (This InterfaceImport) DiffInterface(oldInterfaceList []interfaces.InterfaceTable) (code int, err error) {
	tx := global.DBMysql.Begin()
	for _, newInterface := range This.InterfaceList {
		temp := true
		for _, oldInterface := range oldInterfaceList {
			// 2.修改已有的
			if newInterface.Path == oldInterface.Path && newInterface.Method == oldInterface.Method {
				temp = false
				if code, err = This.UpdateInterface(tx, newInterface, oldInterface); err != nil {
					tx.Rollback()
					return code, err
				}
			}
		}
		// 1.增加没有的
		if temp {
			if code, err = This.InsertInterface(tx, newInterface); err != nil {
				tx.Rollback()
				return code, err
			}
		}
	}
	for _, oldInterface := range oldInterfaceList {
		temp := true
		for _, newInterface := range This.InterfaceList {
			if newInterface.Path == oldInterface.Path && newInterface.Method == oldInterface.Method {
				temp = false
			}
		}
		if temp {
			// 3.删除多余的
			if code, err = This.DelInterface(tx, oldInterface); err != nil {
				tx.Rollback()
				return code, err
			}
		}
	}
	tx.Commit()
	return
}

func (This InterfaceImport) InsertInterface(tx *gorm.DB, newInterface InterfaceInfo) (code int, err error) {
	interfaceT := interfaces.InterfaceTable{
		InterfaceName: newInterface.InterfaceName,
		InterfaceDesc: newInterface.InterfaceDesc,
		Path:          newInterface.Path,
		Method:        newInterface.Method,
		ServiceId:     This.ServiceId,
	}
	if err = interfaceT.Insert(tx); err != nil {
		logger.Error("接口新增失败")
		return 2151, errors.New("")
	}
	return
}

func (This InterfaceImport) UpdateInterface(tx *gorm.DB, newInterface InterfaceInfo, oldInterface interfaces.InterfaceTable) (code int, err error) {
	interfaceT := interfaces.InterfaceTable{InterfaceId: oldInterface.InterfaceId}
	upMap := map[string]interface{}{"name": newInterface.InterfaceName}
	if newInterface.InterfaceDesc != "" {
		upMap["desc"] = newInterface.InterfaceDesc
	}
	if err = interfaceT.Update(tx, upMap); err != nil {
		logger.Error("接口修改失败")
		return 2154, errors.New("")
	}
	return
}

func (This InterfaceImport) DelInterface(tx *gorm.DB, oldInterface interfaces.InterfaceTable) (code int, err error) {
	if This.IsDisassociate == 1 {
		// 3.1删除casbin关系
		e := mycasbin.Casbin()
		_, err = e.RemoveFilteredPolicy(1, oldInterface.Path, oldInterface.Method)
		if err != nil {
			logger.Error("角色与casbin联系删除失败")
			return 2115, errors.New("")
		}
		// 3.2删除角色与接口关系
		roleInterfaceT := roles.RoleInterfaceTable{}
		if err = roleInterfaceT.DeleteByInterfaceIds(tx, []int{oldInterface.InterfaceId}); err != nil {
			logger.Error("角色与casbin联系删除失败")
			return 2113, errors.New("")
		}
	}
	// 3.3删除接口
	interfaceT := interfaces.InterfaceTable{InterfaceId: oldInterface.InterfaceId}
	if err = interfaceT.Delete(tx); err != nil {
		logger.Error("接口删除失败", oldInterface.InterfaceId)
		return 2153, errors.New("")
	}
	return
}

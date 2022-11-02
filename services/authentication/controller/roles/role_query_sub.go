package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type RoleSubQuery struct {
	RoleId int `json:"role_id"       binding:"required"`
}

type RoleSubQueryRes struct {
	ServiceId   int       `json:"service_id"`
	ServiceName string    `json:"service_name"`
	RoleSubList []RoleSub `json:"role_list"`
}

type RoleSub struct {
	RoleId     int    `json:"role_id"`
	RoleName   string `json:"role_name"`
	Status     int    `json:"status"`
	RoleStatus int    `json:"role_status"`
}

func (This RoleSubQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色查询_小角色 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.查询所有
	res, code, err := This.GetRoles()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.分配角色
	resEd, code, err := This.GetRolesEd(res)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ(resEd)
}

func (This RoleSubQuery) JudgeInfo(userId int) (code int, err error) {
	// 1.该用户是否是超级管理员
	if err = model.JudgeAuthManager(userId); err != nil {
		logger.Error("权限不足")
		return 2171, err
	}
	return
}

func (This RoleSubQuery) GetRoles() (res []RoleSubQueryRes, code int, err error) {
	// 1.获取所有服务id
	serviceT := services.ServiceTable{}
	serviceList, err := serviceT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Info("服务查询失败")
		return res, 2101, errors.New("")
	}
	// 2.获取角色
	roleT := roles.RoleTable{}
	roleList, err := roleT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Info("角色查询失败")
		return res, 2112, errors.New("")
	}
	// 3.获取角色绑定数
	roleNumMap := map[int]int{}
	roleSubT := roles.RoleSubTable{}
	roleSubList, err := roleSubT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Info("角色与大角色联系查询失败")
		return res, 2118, errors.New("")
	}
	for _, roleSub := range roleSubList {
		num, ok := roleNumMap[roleSub.SubRoleId]
		if ok {
			roleNumMap[roleSub.SubRoleId] = num + 1
		} else {
			roleNumMap[roleSub.SubRoleId] = 1
		}
	}
	// 4.获取所有role
	for _, service := range serviceList {
		// 4.1除Auth外
		if service.ServiceId == global.DefaultServiceId {
			continue
		}
		// 4.2除csr其他服务
		if utils.IsContainsInt(global.CsrServiceIds, service.ServiceId) {
			continue
		}
		roleListRes := RoleSubQueryRes{ServiceId: service.ServiceId, ServiceName: service.ServiceDesc}
		for _, role := range roleList {
			if service.ServiceId == role.ServiceId {
				roleSub := RoleSub{RoleId: role.RoleId, RoleName: role.RoleName, RoleStatus: 1}
				// 4.3判断角色数是否可以绑定
				if role.RoleNum != -1 {
					num, ok := roleNumMap[role.RoleId]
					if ok {
						if role.RoleNum <= num {
							roleSub.RoleStatus = 0
						}
					} else {
						if role.RoleNum == 0 {
							roleSub.RoleStatus = 0
						}
					}
				}
				roleListRes.RoleSubList = append(roleListRes.RoleSubList, roleSub)
			}
		}
		if len(roleListRes.RoleSubList) == 0 {
			roleListRes.RoleSubList = []RoleSub{}
		}
		res = append(res, roleListRes)
	}
	return
}

func (This RoleSubQuery) GetRolesEd(res []RoleSubQueryRes) (resEd []RoleSubQueryRes, code int, err error) {
	// 1.获取角色所绑定小角色
	roleSubT := roles.RoleSubTable{}
	roleSubList, err := roleSubT.QueryListByFilter(global.DBMysql, map[string]interface{}{"role_id": This.RoleId})
	if err != nil {
		logger.Info("角色与大角色联系查询失败")
		return res, 2118, errors.New(fmt.Sprint("RoleId:", This.RoleId))
	}
	for i, roleSubQueryRes := range res {
		for j, roleSubQuery := range roleSubQueryRes.RoleSubList {
			for _, roleSub := range roleSubList {
				if roleSub.SubRoleId == roleSubQuery.RoleId {
					res[i].RoleSubList[j].Status = 1
					res[i].RoleSubList[j].RoleStatus = 1
					break
				}
			}
		}
	}
	// 2.大屏服务需在最下
	screenMax := RoleSubQueryRes{}
	for _, roleSubQueryRes := range res {
		if roleSubQueryRes.ServiceId == 3 {
			screenMax = roleSubQueryRes
			continue
		}
		resEd = append(resEd, roleSubQueryRes)
	}
	resEd = append(resEd, screenMax)
	return
}

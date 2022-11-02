package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserTokenQuery struct {
	ServiceName string `json:"service_name" binding:"required"`
}

type UserTokenQueryRes struct {
	UserId             int    `json:"user_id"`
	Account            string `json:"account"`
	UserName           string `json:"user_name"`
	RoleId             int    `json:"role_id"`
	CreateTime         int64  `json:"create_time"`
	Manager            int    `json:"manager"`
	DepartmentId       int    `json:"department_id"`
	ManageDepartmentId []int  `json:"manage_department_id"`
}

func (This UserTokenQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("用户token查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.先获取信息
	userT := users.UserTable{}
	user, err := userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": userId})
	if err != nil {
		logger.Error("用户查询失败")
		ico.Err(2144, "")
	}
	// 2.查询部角色信息
	roleT := users.UserRoleTable{}
	role, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"user_id": userId})
	if err != nil {
		logger.Error("角色查询失败")
		ico.Err(2144, "")
	}
	// 3.查询部门信息
	departmentT := departments.DepartmentTab{}
	departmentList, err := departmentT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error(err)
		ico.Err(2148, "")
	}
	manageDepartmentIds := []int{}
	for _, department := range departmentList {
		if department.ChargerId == userId {
			manageDepartmentIds = append(manageDepartmentIds, department.Id)
		}
	}
	res := UserTokenQueryRes{
		UserId: userId, Account: user.Account, UserName: user.UserName, CreateTime: user.CreateTime,
		DepartmentId: user.DepartmentId, ManageDepartmentId: manageDepartmentIds, RoleId: role.RoleId,
	}
	// 4.1 查询所有
	if This.ServiceName == "-1" {
		resMap, code, err := This.GetUsers(userId, res)
		if err != nil {
			return ico.Err(code, err.Error())
		}
		return ico.Succ(resMap)
	}
	// 4.2 查询单个 该用户是否是该服务管理员
	if err = model.JudgeManager(userId, This.ServiceName); err == nil {
		res.Manager = 1
	}
	return ico.Succ(res)
}

func (This UserTokenQuery) GetUsers(userId int, res UserTokenQueryRes) (resMap map[string]UserTokenQueryRes, code int, err error) {
	// 1.获取所有服务id
	serviceT := services.ServiceTable{}
	serviceList, err := serviceT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Info("服务查询失败")
		return resMap, 2101, errors.New("")
	}
	// 2.配置管理员
	resMap = make(map[string]UserTokenQueryRes)
	for _, service := range serviceList {
		res.Manager = 0
		// 2.1除csr其他服务
		if utils.IsContainsInt(global.CsrServiceIds, service.ServiceId) {
			continue
		}
		// 2.2该用户是否是该服务管理员
		if err = model.JudgeManager(userId, service.ServiceName); err == nil {
			res.Manager = 1
		}
		resMap[service.ServiceName] = res
	}
	err = nil
	return
}

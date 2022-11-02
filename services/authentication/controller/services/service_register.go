package services

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ServiceRegister struct {
	ServiceName string `json:"service_name" binding:"required"`
	IndexUrl    string `json:"index_url"`
	ServiceDesc string `json:"service_desc" binding:"required"`
}

type ServiceRegisterRes struct {
	ServiceId int `json:"service_id"`
}

func (This ServiceRegister) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	logger.Infof("服务注册")
	// 1.校验
	if code, err := This.JudgeInfo(); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.服务是否存在
	serviceT := services.ServiceTable{}
	serviceList, err := serviceT.QueryListByFilter(global.DBMysql, map[string]interface{}{"name": This.ServiceName})
	if err != nil {
		return ico.Err(2101, "")
	}
	res := ServiceRegisterRes{}
	// 非首次注册
	if len(serviceList) == 1 {
		if code, err := This.UpdateDefault(serviceList[0].ServiceId); err != nil {
			return ico.Err(code, err.Error())
		}
		res.ServiceId = serviceList[0].ServiceId
		return ico.Succ(res)
	}
	// 首次注册
	tx := global.DBMysql.Begin()
	serviceId, code, err := This.InsertDefault(tx)
	if err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	res.ServiceId = serviceId
	tx.Commit()
	return ico.Succ(res)
}

func (This ServiceRegister) JudgeInfo() (code int, err error) {
	// 需以 -service结尾
	if !strings.HasSuffix(This.ServiceName, "-service") || This.ServiceName == "-service" {
		logger.Info("服务名称格式异常")
		return 2103, errors.New("")
	}
	return
}

func (This ServiceRegister) UpdateDefault(serviceId int) (code int, err error) {
	// 覆盖其他
	serviceT := services.ServiceTable{ServiceId: serviceId, ServiceDesc: This.ServiceDesc, IndexUrl: This.IndexUrl}
	if err = serviceT.UpdateByStruct(global.DBMysql); err != nil {
		logger.Info("服务修改失败")
		return 2103, errors.New("")
	}
	return
}

func (This ServiceRegister) InsertDefault(tx *gorm.DB) (serviceId, code int, err error) {
	// 首次注册
	serviceT := services.ServiceTable{
		ServiceName: This.ServiceName,
		IndexUrl:    This.IndexUrl,
		ServiceDesc: This.ServiceDesc,
	}
	// 1.新增服务
	if err = serviceT.Insert(tx); err != nil {
		logger.Info("服务新增失败")
		return serviceT.ServiceId, 2102, errors.New("")
	}
	// 2.新增内置角色
	roleT := roles.RoleTable{
		RoleName:   fmt.Sprintf("管理员"),
		RoleDesc:   fmt.Sprintf("%s 内置管理员", This.ServiceName),
		ServiceId:  serviceT.ServiceId,
		Manager:    global.DefaultManagerInt,
		UserNum:    global.DefaultLimitNum,
		RoleNum:    global.DefaultLimitNum,
		CreateTime: time.Now().Unix(),
	}
	if err = roleT.Insert(tx); err != nil {
		logger.Info("角色新增失败")
		return serviceT.ServiceId, 2111, errors.New("")
	}
	// 3.新增内置账户
	//password, _ := utils.HashPassword(fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix))
	//userT := users.UserTable{
	//	Account:    fmt.Sprintf("%s_%s", This.ServiceName, global.DefaultSuffix),
	//	UserName:   fmt.Sprintf("%s_内置管理员账户", This.ServiceName),
	//	UserDesc:   fmt.Sprintf("%s_内置管理员账户", This.ServiceName),
	//	Password:   password,
	//	Default:    global.DefaultUserInt,
	//	CreateTime: time.Now().Unix(),
	//}
	//if err = userT.Insert(tx); err != nil {
	//	logger.Info("服务新增失败")
	//	return serviceT.ServiceId, 2131, errors.New("用户新增失败")
	//}
	//// 4.为内置账号添加角色
	//userRoleT := users.UserRoleTable{
	//	UserId: userT.UserId,
	//	RoleId: roleT.RoleId,
	//}
	//if err = userRoleT.Insert(tx); err != nil {
	//	logger.Info("服务新增失败")
	//	return serviceT.ServiceId, 2140, errors.New("用户与角色联系修改失败")
	//}
	serviceId = serviceT.ServiceId
	return
}

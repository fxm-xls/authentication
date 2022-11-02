package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/config"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserJFAdd struct {
	SysId       string        `json:"sysId"         binding:"required"`
	ReqTime     string        `json:"reqTime"`
	AddCount    string        `json:"addCount"      binding:"required"`
	DelCount    string        `json:"delCount"      binding:"required"`
	AddOperList []AddOperInfo `json:"addOperList"`
	DelOperList []string      `json:"delOperList"`
}

type AddOperInfo struct {
	OperId     string `json:"operId"`
	OperName   string `json:"operName"`
	Msisdn     string `json:"msisdn"`
	Email      string `json:"email"`
	Sex        string `json:"sex"`
	RegionName string `json:"regionName"`
	OrgName    string `json:"orgName"`
}

func (This UserJFAdd) DoHandle(c *gin.Context) *ico.JFResult {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.JFErr(2099, "")
	}
	logger.Infof("用户新增_经分")
	logger.Info("增加用户列表：", This.AddCount, This.AddOperList)
	logger.Info("删除用户列表：", This.DelCount, This.DelOperList)
	// 1.校验
	ip := c.Request.Header.Get("sourceIp")
	if code, err := This.JudgeInfo(ip); err != nil {
		return ico.JFErr(code, err.Error())
	}
	var tx = global.DBMysql.Begin()
	// 2.新增用户
	for _, operInfo := range This.AddOperList {
		if code, err := This.InsertUser(tx, operInfo); err != nil {
			tx.Rollback()
			return ico.JFErr(code, err.Error())
		}
	}
	// 3.删除用户
	for _, operId := range This.DelOperList {
		if code, err := This.DeleteUser(tx, operId); err != nil {
			tx.Rollback()
			return ico.JFErr(code, err.Error())
		}
	}
	tx.Commit()
	return ico.JFSucc()
}

func (This UserJFAdd) JudgeInfo(ipNew string) (code int, err error) {
	// 1.检测是否是经分token
	if This.SysId != "H00000000027" {
		logger.Info(This.SysId, "服务不存在")
		return 2104, errors.New(This.SysId)
	}
	// 2.检测ip
	temp := false
	for _, ip := range config.JFConfig.IPS {
		if ip == ipNew {
			temp = true
		}
	}
	if !temp {
		logger.Info("权限不足:该ip未授权-" + ipNew)
		return 2171, errors.New("该ip未授权")
	}
	// 3.检测增加，删除数量
	addCount, err := strconv.Atoi(This.AddCount)
	if err != nil {
		logger.Info(This.AddCount, "参数异常:增加数非整形")
		return 2099, errors.New("增加数非整形")
	}
	delCount, err := strconv.Atoi(This.DelCount)
	if err != nil {
		logger.Info(This.DelCount, "参数异常:删除数非整形")
		return 2099, errors.New("删除数非整形")
	}
	if addCount != len(This.AddOperList) || delCount != len(This.DelOperList) {
		logger.Info("参数异常:数量不一致")
		return 2099, errors.New("数量不一致")
	}
	return
}

func (This UserJFAdd) InsertUser(tx *gorm.DB, operInfo AddOperInfo) (code int, err error) {
	// hash密码
	password, err := utils.HashPassword(operInfo.OperId)
	if err != nil {
		logger.Info(operInfo.OperId, "用户新增失败")
		return 2132, errors.New("")
	}
	userT := users.UserTable{
		Account:      operInfo.OperId,
		UserName:     operInfo.OperName,
		Password:     password,
		MobilePhone:  operInfo.Msisdn,
		Email:        operInfo.Email,
		Sex:          operInfo.Sex,
		City:         operInfo.RegionName,
		Department:   operInfo.OrgName,
		DepartmentId: config.JFConfig.DeptId,
		CreateTime:   time.Now().Unix(),
		Default:      global.DefaultJFUserInt,
	}
	// 0.判断用户是否存在
	userList, _ := userT.QueryListByFilter(global.DBMysql, map[string]interface{}{"account": operInfo.OperId})
	if len(userList) != 0 {
		logger.Info("账户已存在:" + operInfo.OperId)
		return 200, nil
	}
	// 0.1 判断部门是否存在
	deptT := departments.DepartmentTab{}
	deptInfo, err := deptT.QueryByFilter(global.DBMysql, map[string]interface{}{"dept_name": operInfo.OrgName})
	if err != nil {
		logger.Info("部门不存在: " + operInfo.OrgName)
		err = nil
	} else {
		userT.DepartmentId = deptInfo.Id
	}
	// 1.增加用户信息
	if err = userT.Insert(tx); err != nil {
		logger.Info("用户新增失败")
		return 2132, errors.New("")
	}
	// 2.增加用户角色
	userRoleT := users.UserRoleTable{UserId: userT.UserId, RoleId: config.JFConfig.RoleId}
	if err = userRoleT.Insert(tx); err != nil {
		logger.Info("用户与角色联系修改失败")
		return 2140, errors.New("")
	}
	return
}

func (This UserJFAdd) DeleteUser(tx *gorm.DB, account string) (code int, err error) {
	// 0.判断用户是否存在
	userT := users.UserTable{}
	user, err := userT.QueryByFilter(tx, map[string]interface{}{"account": account})
	if err != nil || user.Default != global.DefaultJFUserInt {
		logger.Info("账户不存在:" + account)
		return 200, nil
	}
	// 1.删除用户数据
	userEntityT := users.UserEntityTable{}
	if err = userEntityT.DeleteByUserIds(tx, []int{user.UserId}); err != nil {
		logger.Info("用户与数据联系删除失败")
		return 2137, errors.New("")
	}
	// 2.删除用户角色
	userRoleT := users.UserRoleTable{}
	if err = userRoleT.DeleteByUserIds(tx, []int{user.UserId}); err != nil {
		logger.Info("用户与角色联系删除失败")
		return 2136, errors.New("")
	}
	// 3.删除用户信息
	if err = userT.DeleteByFilter(tx, map[string]interface{}{"account": account}); err != nil {
		logger.Info("用户删除失败")
		return 2138, errors.New("")
	}
	// 4.删除对应token
	if err = users.DelToken(tx, user.UserId); err != nil {
		logger.Error("用户token删除失败")
		return 2146, errors.New("")
	}
	return
}

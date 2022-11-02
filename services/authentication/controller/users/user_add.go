package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserAdd struct {
	Account            string `json:"account"              binding:"required"`
	UserName           string `json:"user_name"            binding:"required"`
	Password           string `json:"password"             binding:"required"`
	MobilePhone        string `json:"mobilephone"`
	RoleId             int    `json:"role_id"              binding:"required"`
	DepartmentId       int    `json:"department_id"        binding:"required"`
	ManageDepartmentId []int  `json:"manage_department_id"`
	Message            string `json:"message"              binding:"required"`
}

type UserAddRes struct {
	UserId int `json:"user_id"`
}

func (This UserAdd) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	deptId := c.GetInt("dept_id")
	logger.Infof("用户新增 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	if err := model.JudgeAuthManager(userId); err == nil {
		// 2.超级管理员直接操作
		// 2.1 新增用户
		var tx = global.DBMysql.Begin()
		newUserId, code, err := This.InsertUser(tx)
		if err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		// 2.2 增加用户角色联系
		if code, err = This.InsertUserRole(tx, newUserId); err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		// 2.3 增加用户管理部门
		if code, err = This.InsertUserDepartment(tx, newUserId); err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		tx.Commit()
		return ico.Succ(UserAddRes{UserId: newUserId})
	} else {
		// 3.管理员需要提交审批
		if code, err := This.InsertUserApproval(userId, deptId); err != nil {
			return ico.Err(code, err.Error())
		}
		return ico.Succ("")
	}
}

func (This UserAdd) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if code, err = public.JudgeManager(userId); err != nil {
		logger.Error("权限不足")
		return code, errors.New("")
	}
	// 2.账户名称不能为空字符
	if strings.TrimSpace(This.UserName) == "" || strings.TrimSpace(This.Account) == "" {
		logger.Error("账户名称不能为空字符")
		return 2157, errors.New("名称或描述不能为空字符")
	}
	// 2.检测账户是否存在
	userT := users.UserTable{}
	if _, err = userT.QueryByFilter(global.DBMysql, map[string]interface{}{"account": This.Account}); err == nil {
		logger.Error("账户已存在")
		return 2131, errors.New("")
	}
	// 3.角色已存在且是Auth服务
	roleT := roles.RoleTable{}
	role, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId, "service_id": global.DefaultServiceId})
	if err != nil {
		logger.Error("角色不存在")
		return 2119, errors.New("")
	}
	// 4.判断该角色 用户人数是否限制
	if role.UserNum != -1 {
		userNumMax := role.UserNum
		userRoleT := users.UserRoleTable{}
		userRoleList, err := userRoleT.QueryListByFilter(global.DBMysql, map[string]interface{}{"role_id": This.RoleId})
		if err != nil {
			logger.Error("角色查询失败")
			return 2112, errors.New("")
		}
		if userNumMax <= len(userRoleList) {
			logger.Error("角色限制用户数超出")
			return 2124, errors.New("")
		}
	}
	// 5.判断部门是否存在
	departmentT := departments.DepartmentTab{}
	whereMap := map[string]interface{}{"id": This.DepartmentId}
	_, err = departmentT.QueryByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Error("查询部门数据为空")
		return 2207, errors.New("")
	}
	// 6.判断所管部门是否已有管理员
	for _, departmentId := range This.ManageDepartmentId {
		whereMap = map[string]interface{}{"id": departmentId}
		departmentInfo, err := departmentT.QueryByFilter(global.DBMysql, whereMap)
		if err != nil {
			logger.Error("查询部门数据为空")
			return 2207, errors.New("")
		}
		if departmentInfo.ChargerId != 0 {
			logger.Error(departmentInfo.Id, departmentInfo.DeptName, "所选部门已有管理员")
			return 2207, errors.New(departmentInfo.DeptName + " 已有管理员")
		}
	}
	return
}

func (This UserAdd) InsertUser(tx *gorm.DB) (userId, code int, err error) {
	// hash密码
	password, err := utils.HashPassword(This.Password)
	if err != nil {
		logger.Error("用户新增失败")
		return userId, 2132, errors.New("")
	}
	userT := users.UserTable{
		Account:      This.Account,
		UserName:     This.UserName,
		Password:     password,
		DepartmentId: This.DepartmentId,
		MobilePhone:  This.MobilePhone,
		CreateTime:   time.Now().Unix(),
	}
	if err = userT.Insert(tx); err != nil {
		logger.Error("用户新增失败")
		return userId, 2132, errors.New("")
	}
	userId = userT.UserId
	return
}

func (This UserAdd) InsertUserRole(tx *gorm.DB, userId int) (code int, err error) {
	userRoleT := users.UserRoleTable{UserId: userId, RoleId: This.RoleId}
	if err = userRoleT.Insert(tx); err != nil {
		logger.Error("用户与角色联系修改失败")
		return 2140, errors.New("")
	}
	return
}

func (This UserAdd) InsertUserDepartment(tx *gorm.DB, userId int) (code int, err error) {
	for _, departmentId := range This.ManageDepartmentId {
		departmentT := departments.DepartmentTab{Id: departmentId}
		if err = departmentT.Update(tx, map[string]interface{}{"charger_id": userId}); err != nil {
			logger.Error("添加部门管理员失败")
			return 2149, errors.New("")
		}
	}
	return
}

type UserAddApproval struct {
	Account            string `json:"account"`
	UserName           string `json:"user_name"`
	Pwd                string `json:"pwd"`
	DeptId             int    `json:"dept_id"`
	RoleId             int    `json:"role_id"`
	ManageDepartmentId []int  `json:"manage_dept_id"`
}

func (This UserAdd) InsertUserApproval(userId, deptId int) (code int, err error) {
	userAddApproval := UserAddApproval{
		Account:            This.Account,
		UserName:           This.UserName,
		DeptId:             This.DepartmentId,
		Pwd:                This.Password,
		RoleId:             This.RoleId,
		ManageDepartmentId: This.ManageDepartmentId,
	}
	contentJson, _ := json.Marshal(&userAddApproval)
	approvalT := approvals.ApprovalTab{
		ApplicantId: userId,
		DeptId:      deptId,
		Type:        1,
		Content:     contentJson,
		Status:      1,
		SubmitMsg:   This.Message,
		CreateTime:  time.Now().Unix(),
	}
	if err = approvalT.Insert(global.DBMysql); err != nil {
		logger.Error("添加用户失败-添加审批失败")
		return 2150, errors.New("")
	}
	return
}

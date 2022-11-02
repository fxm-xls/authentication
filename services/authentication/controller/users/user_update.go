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
	"time"
)

type UserUpdate struct {
	UserId                int    `json:"user_id"       binding:"required"`
	UserName              string `json:"user_name"`
	Password              string `json:"password"`
	ManageDepartmentId    []int  `json:"manage_department_id"`
	RoleId                int    `json:"role_id"`
	Message               string `json:"message"`
	OldManageDepartmentId []int
}

func (This UserUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	deptId := c.GetInt("dept_id")
	logger.Infof("用户修改 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		// 1.1普通用户修改
		if This.UserId == userId {
			code, err = This.UpdateUser(global.DBMysql)
			if err != nil {
				return ico.Err(code, err.Error())
			}
			return ico.Succ("修改用户成功")
		}
		return ico.Err(code, err.Error())
	}
	if err := model.JudgeAuthManager(userId); err == nil {
		// 2.超级管理员直接操作
		// 2.1 修改用户
		var tx = global.DBMysql.Begin()
		if code, err := This.UpdateUser(tx); err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		// 2.2 修改用户角色联系
		if code, err := This.UpdateUserRole(tx); err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		// 2.3 增加用户管理部门
		if code, err := This.UpdateUserDepartment(tx); err != nil {
			tx.Rollback()
			return ico.Err(code, err.Error())
		}
		tx.Commit()
	} else {
		// 3.管理员需要提交审批
		if code, err := This.UpdateUserApproval(userId, deptId); err != nil {
			return ico.Err(code, err.Error())
		}
		return ico.Succ("")
	}
	return ico.Succ("修改用户成功")
}

func (This UserUpdate) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if code, err = public.JudgeManager(userId); err != nil {
		logger.Error("权限不足")
		return code, errors.New("")
	}
	// 2.检测账户是否存在
	userT := users.UserTable{}
	user, err := userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.UserId})
	if err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	// 3.检测账户是否是内置
	if user.Default == global.DefaultUserInt {
		logger.Error("内置用户")
		return 2144, errors.New("内置用户-不可修改")
	}
	// 4.检测角色
	if This.RoleId != 0 {
		// 4.1 角色已存在且是Auth服务
		roleT := roles.RoleTable{}
		role, err := roleT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.RoleId, "service_id": global.DefaultServiceId})
		if err != nil {
			logger.Error("角色不存在")
			return 2119, errors.New("")
		}
		// 4.2 判断该角色人数是否限制
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
	}
	// 5.判断部门是否存在 判断所管部门是否已有管理员且管理员为其他
	departmentT := departments.DepartmentTab{}
	for _, departmentId := range This.ManageDepartmentId {
		whereMap := map[string]interface{}{"id": departmentId}
		departmentInfo, err := departmentT.QueryByFilter(global.DBMysql, whereMap)
		if err != nil {
			logger.Error("查询部门数据为空")
			return 2207, errors.New("")
		}
		if departmentInfo.ChargerId != 0 && departmentInfo.ChargerId != This.UserId {
			logger.Error(departmentInfo.Id, departmentInfo.DeptName, " 所选部门已有其他管理员")
			return 2207, errors.New(departmentInfo.DeptName + " 已有管理员")
		}
	}
	return
}

func (This UserUpdate) UpdateUser(tx *gorm.DB) (code int, err error) {
	userT := users.UserTable{
		UserId:   This.UserId,
		UserName: This.UserName,
	}
	if This.Password != "" {
		// hash密码
		password, err := utils.HashPassword(This.Password)
		if err != nil {
			logger.Error("用户信息修改失败")
			return 2139, errors.New("")
		}
		userT.Password = password
	}
	if err = userT.UpdateByStruct(tx); err != nil {
		logger.Error("用户信息修改失败")
		return 2139, errors.New("")
	}
	return
}

func (This UserUpdate) UpdateUserRole(tx *gorm.DB) (code int, err error) {
	if This.RoleId == 0 {
		return
	}
	userRoleT := users.UserRoleTable{UserId: This.UserId, RoleId: This.RoleId}
	// 1.先删除
	if err = userRoleT.DeleteByUserIds(tx, []int{This.UserId}); err != nil {
		logger.Error("用户与角色联系删除失败")
		return 2136, errors.New("")
	}
	// 2.再新增
	if err = userRoleT.Insert(tx); err != nil {
		logger.Error("用户与角色联系修改失败")
		return 2140, errors.New("")
	}
	return
}

func (This UserUpdate) UpdateUserDepartment(tx *gorm.DB) (code int, err error) {
	departmentT := departments.DepartmentTab{}
	// 1.获取旧管理部门id列表
	departmentList, err := departmentT.QueryListByFilter(tx, map[string]interface{}{"charger_id": This.UserId})
	if err != nil {
		logger.Error("查询部门数据失败")
		return 2202, errors.New("")
	}
	for _, department := range departmentList {
		This.OldManageDepartmentId = append(This.OldManageDepartmentId, department.Id)
	}
	// 2.删除管理部门id
	for _, departmentOldId := range This.OldManageDepartmentId {
		if !utils.IsContainsInt(This.ManageDepartmentId, departmentOldId) {
			departmentT = departments.DepartmentTab{Id: departmentOldId}
			if err = departmentT.Update(tx, map[string]interface{}{"charger_id": nil}); err != nil {
				logger.Error("删除部门管理员失败")
				return 2149, errors.New("")
			}
		}
	}
	// 3.新增管理部门id
	for _, departmentId := range This.ManageDepartmentId {
		if !utils.IsContainsInt(This.OldManageDepartmentId, departmentId) {
			departmentT = departments.DepartmentTab{Id: departmentId}
			if err = departmentT.Update(tx, map[string]interface{}{"charger_id": This.UserId}); err != nil {
				logger.Error("添加部门管理员失败")
				return 2149, errors.New("")
			}
		}
	}
	return
}

func (This UserUpdate) UpdateUserApproval(userId, deptId int) (code int, err error) {
	// 获取用户旧信息
	userV := users.UserView{}
	userInfo, err := userV.QueryByFilter(global.DBMysql, map[string]interface{}{"user_id": This.UserId})
	if err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	content := map[string]interface{}{
		"user_id":            This.UserId,
		"user_name":          userInfo.UserName,
		"user_name_new":      This.UserName,
		"role_id":            userInfo.RoleId,
		"role_id_new":        This.RoleId,
		"manage_dept_id":     This.OldManageDepartmentId,
		"manage_dept_id_new": This.ManageDepartmentId,
	}
	contentJson, _ := json.Marshal(&content)
	approvalT := approvals.ApprovalTab{
		ApplicantId: userId,
		DeptId:      deptId,
		Type:        2,
		Content:     contentJson,
		Status:      1,
		SubmitMsg:   This.Message,
		CreateTime:  time.Now().Unix(),
	}
	if err = approvalT.Insert(global.DBMysql); err != nil {
		logger.Error("修改用户失败-添加审批失败")
		return 2150, errors.New("")
	}
	return
}

package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDelete struct {
	UserId int `json:"user_id" binding:"required"`
}

func (This UserDelete) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("删除用户 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.删除
	var tx = global.DBMysql.Begin()
	if code, err := This.DeleteUser(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("删除用户成功")
}

func (This UserDelete) JudgeInfo(userId int) (code int, err error) {
	// 1.管理员检测
	if code, err = public.JudgeManager(userId); err != nil {
		logger.Error("权限不足")
		return code, errors.New("")
	}
	// 2.自杀检测
	if This.UserId == userId {
		logger.Error("自杀失败")
		return 2145, errors.New("")
	}
	// 3.检测账户是否存在
	userT := users.UserTable{}
	user, err := userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.UserId})
	if err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	// 4.检测账户是否是内置
	if user.Default == global.DefaultUserInt {
		logger.Error("用户不存在-内置用户不可删除")
		return 2144, errors.New("内置用户不可删除")
	}
	// 5.检测账户是否有所管部门
	if err = model.JudgeDepartmentManager(This.UserId); err == nil {
		logger.Error("账户有所管部门")
		return 2138, errors.New("账户有所管部门")
	}
	err = nil
	return
}

func (This UserDelete) DeleteUser(tx *gorm.DB) (code int, err error) {
	// 1.删除用户信息
	userT := users.UserTable{UserId: This.UserId}
	if err = userT.Delete(tx); err != nil {
		logger.Error("用户删除失败")
		return 2138, errors.New("")
	}
	// 2.删除对应数据联系
	userEntityT := users.UserEntityTable{}
	if err = userEntityT.DeleteByUserIds(tx, []int{This.UserId}); err != nil {
		logger.Error("用户与数据联系删除失败")
		return 2137, errors.New("")
	}
	// 3.删除对应角色联系
	userRoleT := users.UserRoleTable{}
	if err = userRoleT.DeleteByUserIds(tx, []int{This.UserId}); err != nil {
		logger.Error("用户与角色联系删除失败")
		return 2136, errors.New("")
	}
	// 4.删除对应部门联系
	if code, err = This.DeleteUserDepartment(tx); err != nil {
		logger.Error("用户与部门联系删除失败")
		return 2149, errors.New("")
	}
	// 5.删除对应token
	if err = users.DelToken(tx, This.UserId); err != nil {
		logger.Error("用户token删除失败")
		return 2146, errors.New("")
	}
	return
}

func (This UserDelete) DeleteUserDepartment(tx *gorm.DB) (code int, err error) {
	departmentT := departments.DepartmentTab{}
	departmentList, err := departmentT.QueryListByFilter(tx, map[string]interface{}{"charger_id": This.UserId})
	if err != nil {
		logger.Error("查询部门数据失败")
		return 2202, errors.New("")
	}
	for _, department := range departmentList {
		departmentT = departments.DepartmentTab{Id: department.Id}
		if err = departmentT.Update(tx, map[string]interface{}{"charger_id": nil}); err != nil {
			logger.Error("删除部门管理员失败")
			return 2149, errors.New("删除部门管理员失败")
		}
	}
	return
}

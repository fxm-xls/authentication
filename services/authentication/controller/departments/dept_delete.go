package departments

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DeptDelete struct {
	DepartmentId int `json:"department_id"`
}

func (d DeptDelete) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-delete department=======================")
	if err := c.ShouldBindJSON(&d); err != nil {
		logger.Error(err.Error())
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	deptId := c.GetInt("dept_id")
	logger.Infof("登录用户基本信息 userId: %d, deptId %d", userId, deptId)
	// 部门管理员访问权限验证
	if code, err := public.JudgeManager(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	tabUser := users.UserTable{}
	dsUser, err := tabUser.QueryByFilter(global.DBMysql, map[string]interface{}{"dept_id": d.DepartmentId})
	logger.Info("[INFO] " + fmt.Sprintf("dsUser: %v", dsUser))
	if err != nil {
		logger.Error(err.Error())
	}
	if dsUser != (users.UserTable{}) {
		// 删除部门有人员不允许删除
		return ico.Err(2209, "", errors.New("delete department personnel cannot delete"))
	}

	tx := global.DBMysql.Begin()
	tabDept := departments.DepartmentTab{Id: d.DepartmentId}
	if err := tabDept.Delete(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
		return ico.Err(2206, "", err.Error())
	}

	tabDeptRole := departments.DeptRoleTab{DeptId: d.DepartmentId}
	if err := tabDeptRole.Delete(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
		return ico.Err(2213, "", err.Error())
	}
	tx.Commit()

	return ico.Succ("部门删除成功")
}

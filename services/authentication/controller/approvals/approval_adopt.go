package approvals

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type ApprovalAdopt struct {
	List []AdoptList `json:"list"`
}

type AdoptList struct {
	ApprovalId int `json:"approval_id"`
}

func (a ApprovalAdopt) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-query approval_adopt=======================")
	if err := c.ShouldBindJSON(&a); err != nil {
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

	for _, v := range a.List {
		// write table
		tabApproval := approvals.ApprovalTab{}
		dsApproval, err := tabApproval.QueryByFilter(global.DBMysql, map[string]interface{}{"id": v.ApprovalId})
		logger.Info("[INFO] " + fmt.Sprintf("dsApproval: %v", dsApproval))
		if err != nil {
			logger.Error(err.Error())
			return ico.Err(2220, "", err.Error())
		}
		if dsApproval.Id == 0 {
			return ico.Err(2215, "", errors.New("no approval record is found"))
		}
		switch dsApproval.Type {
		case 1: // 新增账号
			ua := public.UserAdd{}
			if err := json.Unmarshal(dsApproval.Content, &ua); err != nil {
				logger.Error(err)
			}
			if err := a.doUserAdd(ua); err != nil {
				return ico.Err(2218, "", err.Error())
			}
		case 2: // 修改账号
			uu := public.UserUpdate{}
			if err := json.Unmarshal(dsApproval.Content, &uu); err != nil {
				logger.Error(err)
			}
			if err := a.doUserUpdate(uu); err != nil {
				return ico.Err(2219, "", err.Error())
			}
		case 3: // 新增部门
			da := public.DeptAdd{}
			if err := json.Unmarshal(dsApproval.Content, &da); err != nil {
				logger.Error(err)
			}
			if err := a.doDeptAdd(da); err != nil {
				return ico.Err(2216, "", err.Error())
			}
		case 4: // 修改部门
			du := public.DeptUpdate{}
			if err := json.Unmarshal(dsApproval.Content, &du); err != nil {
				logger.Error(err)
			}
			if err := a.doDeptUpdate(du); err != nil {
				return ico.Err(2217, "", err.Error())
			}
		}
		// update status & message
		tx := global.DBMysql.Begin()
		data := map[string]interface{}{
			"status":        2,
			"approval_time": time.Now().Unix(),
		}
		if err := tabApproval.Update(global.DBMysql, v.ApprovalId, data); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2204, "", err.Error())
		}
	}
	return ico.Succ("审批通过成功")
}

func (a ApprovalAdopt) doUserAdd(us public.UserAdd) error {
	var err error
	tx := global.DBMysql.Begin()
	tabUser := users.UserTable{
		Account:      us.Account,
		UserName:     us.UserName,
		DepartmentId: us.DeptId,
		Password:     us.Pwd,
	}
	if err = tabUser.Insert(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	// TODO manage_dept_id 一个用户管理多个部门
	tabUserRole := users.UserRoleTable{
		UserId: tabUser.UserId,
		RoleId: us.RoleId,
	}
	if err = tabUserRole.Insert(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}
	tx.Commit()

	return err
}

func (a ApprovalAdopt) doUserUpdate(uu public.UserUpdate) error {
	var err error
	tx := global.DBMysql.Begin()
	tabUser := users.UserTable{UserId: uu.UserId}
	data := map[string]interface{}{
		"user_name": uu.UserNameNew,
	}
	if err := tabUser.Update(global.DBMysql, data); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	tabUserRole := users.UserRoleTable{UserId: uu.UserId}
	if err = tabUserRole.Delete(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	tabUserRole = users.UserRoleTable{
		UserId: uu.UserId,
		RoleId: uu.RoleId,
	}
	if err = tabUserRole.Insert(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}
	tx.Commit()

	return err
}

func (a ApprovalAdopt) doDeptAdd(da public.DeptAdd) error {
	var deptIdNew int
	var err error
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryByFilter(global.DBMysql, map[string]interface{}{"dept_name": da.DeptName})
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
	}
	if dsDept != (departments.DepartmentTab{}) {
		logger.Info("[INFO] duplicate departments exist in approval")
		return nil
	}

	tx := global.DBMysql.Begin()
	tabDept = departments.DepartmentTab{
		DeptName:   da.DeptName,
		ParentId:   da.ParentId,
		ChargerId:  da.ChargerId,
		CreateTime: time.Now().Unix(),
		Desc:       "",
		DeptLevel:  da.DeptLevel,
	}
	if deptIdNew, err = tabDept.Insert(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	// 部门创建成功后返回的ID作为DeptId
	for _, v := range da.RoleId {
		tabDeptRole := departments.DeptRoleTab{
			DeptId: deptIdNew,
			RoleId: v,
		}
		if err = tabDeptRole.Insert(global.DBMysql); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
		}
	}
	tx.Commit()

	return err
}

func (a ApprovalAdopt) doDeptUpdate(du public.DeptUpdate) error {
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryByFilter1(global.DBMysql, map[string]interface{}{"dept_name": du.DeptName}, du.DeptId)
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
	}
	if dsDept != (departments.DepartmentTab{}) {
		logger.Info("[INFO] duplicate departments exist in approval")
		return nil
	}

	tx := global.DBMysql.Begin()
	tabDept = departments.DepartmentTab{Id: du.DeptId}
	data := map[string]interface{}{
		"dept_name":  du.DeptNameNew,
		"charger_id": du.ChargerId,
	}
	if err = tabDept.Update(global.DBMysql, data); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	tabDeptRole := departments.DeptRoleTab{DeptId: du.DeptId}
	if err = tabDeptRole.Delete(global.DBMysql); err != nil {
		logger.Error(err.Error())
		tx.Rollback()
	}

	for _, v := range du.RoleId {
		tabDeptRole := departments.DeptRoleTab{
			DeptId: du.DeptId,
			RoleId: v,
		}
		if err = tabDeptRole.Insert(global.DBMysql); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
		}
	}
	tx.Commit()

	return err
}

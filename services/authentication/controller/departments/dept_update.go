package departments

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/model/departments"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type DeptUpdate struct {
	DepartmentId   int    `json:"department_id" binding:"required"`
	DepartmentName string `json:"department_name"`
	ManageId       int    `json:"manage_id"`
	RoleIds        []int  `json:"role_ids"`
	Message        string `json:"message"`
}

func (d DeptUpdate) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-update department=======================")
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
	// 部门名称重名&非空&不含空格
	if d.DepartmentName == "" {
		return ico.Err(2201, "", errors.New("the department name cannot be empty"))
	}
	if strings.ContainsAny(d.DepartmentName, " ") {
		return ico.Err(2221, "", errors.New("the department name cannot contain spaces"))
	}
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryByFilter1(global.DBMysql, map[string]interface{}{"dept_name": d.DepartmentName}, d.DepartmentId)
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
	}
	if dsDept != (departments.DepartmentTab{}) {
		return ico.Err(2203, "", errors.New("duplicate department name"))
	}

	if err := model.JudgeAuthManager(userId); err == nil {
		// 超级管理员不走审批流程
		tx := global.DBMysql.Begin()
		tabDept = departments.DepartmentTab{Id: d.DepartmentId}
		data := map[string]interface{}{
			"dept_name":  d.DepartmentName,
			"charger_id": d.ManageId,
		}
		if err := tabDept.Update(global.DBMysql, data); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2211, "", err.Error())
		}
		tx.Commit()

		tx = global.DBMysql.Begin()
		tabDeptRole := departments.DeptRoleTab{DeptId: d.DepartmentId}
		if err := tabDeptRole.Delete(global.DBMysql); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2213, "", err.Error())
		}
		tx.Commit()

		tx = global.DBMysql.Begin()
		for _, v := range d.RoleIds {
			tabDeptRole := departments.DeptRoleTab{
				DeptId: d.DepartmentId,
				RoleId: v,
			}
			if err := tabDeptRole.Insert(global.DBMysql); err != nil {
				logger.Error(err.Error())
				tx.Rollback()
				return ico.Err(2212, "", err.Error())
			}
		}
		tx.Commit()
	} else {
		obj := public.DeptUpdate{
			DeptId:       d.DepartmentId,
			DeptName:     d.getDeptName(d.DepartmentId),
			DeptNameNew:  d.DepartmentName,
			ChargerId:    d.getChargerId(d.DepartmentId),
			ChargerIdNew: d.ManageId,
			RoleId:       d.getRoleIds(d.DepartmentId),
			RoleIdNew:    d.RoleIds,
		}
		cont, err := json.Marshal(&obj)
		if err != nil {
			return ico.Err(2205, "", err.Error())
		}

		tabApproval := approvals.ApprovalTab{
			ApplicantId: userId,
			DeptId:      deptId,
			Type:        4,
			Content:     cont,
			Status:      1,
			SubmitMsg:   d.Message,
			CreateTime:  time.Now().Unix(),
		}
		tx := global.DBMysql.Begin()
		if err := tabApproval.Insert(global.DBMysql); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2204, "", err.Error())
		}
		tx.Commit()
	}

	return ico.Succ("部门修改成功")
}

func (d DeptUpdate) getDeptName(deptId int) (deptName string) {
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryByFilter(global.DBMysql, map[string]interface{}{"id": deptId})
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	if dsDept.Id == 0 {
		return ""
	}
	deptName = dsDept.DeptName

	return
}

func (d DeptUpdate) getChargerId(deptId int) (chargerId int) {
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryByFilter(global.DBMysql, map[string]interface{}{"id": deptId})
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
		return 0
	}
	if dsDept.Id == 0 {
		return 0
	}
	chargerId = dsDept.ChargerId

	return
}

func (d DeptUpdate) getRoleIds(deptId int) (roleIds []int) {
	tabDeptRole := departments.DeptRoleTab{}
	dsDeptRole, err := tabDeptRole.QueryListByFilter(global.DBMysql, map[string]interface{}{"id": deptId})
	logger.Info("[INFO] " + fmt.Sprintf("dsDeptRole: %v", dsDeptRole))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	if len(dsDeptRole) == 0 {
		return nil
	}
	for _, v := range dsDeptRole {
		roleIds = append(roleIds, v.RoleId)
	}

	return
}

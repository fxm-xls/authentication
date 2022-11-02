package departments

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type DeptAdd struct {
	DepartmentName        string `json:"department_name" binding:"required"`         // 部门名称
	ManageId              int    `json:"manage_id"`                                  // 部门负责人
	RoleIds               []int  `json:"role_ids" binding:"required"`                // 可配置角色
	ParentDepartmentId    int    `json:"parent_department_id"`                       // 上级部门
	Message               string `json:"message" binding:"required"`                 // 申请说明
	ParentDepartmentLevel int    `json:"parent_department_level" binding:"required"` // 上级部门级别
}

func (d DeptAdd) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-add department=======================")
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
	dsDept, err := tabDept.QueryByFilter(global.DBMysql, map[string]interface{}{"dept_name": d.DepartmentName})
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
	}
	if dsDept != (departments.DepartmentTab{}) {
		return ico.Err(2203, "", errors.New("duplicate department name"))
	}
	var deptLevel int
	if d.ParentDepartmentId == 0 {
		deptLevel = 1
	} else if d.ParentDepartmentLevel == 4 {
		return ico.Err(2210, "", errors.New("department max allowed increased to level 4"))
	} else {
		deptLevel = d.ParentDepartmentLevel + 1
	}

	if err := model.JudgeAuthManager(userId); err == nil {
		// 超级管理员不走审批流程
		var deptIdNew int
		var err error
		tx := global.DBMysql.Begin()
		tabDept = departments.DepartmentTab{
			DeptName:   d.DepartmentName,
			ParentId:   d.ParentDepartmentId,
			ChargerId:  d.ManageId,
			CreateTime: time.Now().Unix(),
			Desc:       d.Message,
			DeptLevel:  deptLevel,
		}
		if deptIdNew, err = tabDept.Insert(global.DBMysql); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2211, "", err.Error())
		}
		tx.Commit()

		// 部门创建成功后返回的ID作为DeptId
		tx = global.DBMysql.Begin()
		for _, v := range d.RoleIds {
			tabDeptRole := departments.DeptRoleTab{
				DeptId: deptIdNew,
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
		obj := public.DeptAdd{
			DeptName:  d.DepartmentName,
			ChargerId: d.ManageId,
			RoleId:    d.RoleIds,
			ParentId:  d.ParentDepartmentId,
			DeptLevel: deptLevel,
		}
		cont, err := json.Marshal(&obj)
		if err != nil {
			return ico.Err(2205, "", err.Error())
		}

		tabApproval := approvals.ApprovalTab{
			ApplicantId: userId,
			DeptId:      d.getDeptId(userId),
			Type:        3,
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

	return ico.Succ("部门增加成功")
}

func (d DeptAdd) getDeptId(userId int) (deptId int) {
	tabUser := users.UserTable{}
	dsUser, err := tabUser.QueryByFilter(global.DBMysql, map[string]interface{}{"id": userId})
	logger.Info("[INFO] " + fmt.Sprintf("dsUser: %v", dsUser))
	if err != nil {
		logger.Error(err.Error())
		return 0
	}
	if dsUser.UserId == 0 {
		return 0
	}
	deptId = dsUser.DepartmentId

	return
}

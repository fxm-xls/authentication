package approvals

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/users"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApprovalDetail struct {
	ApprovalId int `json:"approval_id" binding:"required"`
}

type ApprovalDetailResp struct {
	UserInfo       UserInfo       `json:"user_info"`
	DepartmentInfo DepartmentInfo `json:"department_info"`
	ApplicantInfo  ApplicantInfo  `json:"applicant_info"`
	SubmitTime     int64          `json:"submit_time"`
	History        []History      `json:"history"`
	ApprovalType   int            `json:"approval_type"` // 申请类型：1：新增账号，2：修改账号，3：新增部门，4：修改部门
}

type UserInfo struct {
	NewInfo UserInfoNew `json:"new_info"`
	OldInfo UserInfoOld `json:"old_info"`
}

type DepartmentInfo struct {
	NewInfo DepartmentInfoNew `json:"new_info"`
	OldInfo DepartmentInfoOld `json:"old_info"`
}

type ApplicantInfo struct {
	Account          string   `json:"account"`
	UserName         string   `json:"user_name"`
	Department       string   `json:"department"`
	ManageDepartment []string `json:"manage_department"`
	RoleName         string   `json:"role_name"`
}

type History struct {
	Time           int64  `json:"time"`
	UserName       string `json:"user_name"`
	ApprovalStatus int    `json:"approval_status"`
	Message        string `json:"message"`
}

type UserInfoNew struct {
	Account          string   `json:"account"`
	UserName         string   `json:"user_name"`
	Department       string   `json:"department"`
	ManageDepartment []string `json:"manage_department"`
	RoleName         string   `json:"role_name"`
}

type UserInfoOld struct {
	Account          string   `json:"account"`
	UserName         string   `json:"user_name"`
	Department       string   `json:"department"`
	ManageDepartment []string `json:"manage_department"`
	RoleName         string   `json:"role_name"`
}

type DepartmentInfoNew struct {
	DepartmentName   string   `json:"department_name"`
	ManageAccount    string   `json:"manage_account"`
	RoleList         []string `json:"role_list"`
	ParentDepartment string   `json:"parent_department"`
}

type DepartmentInfoOld struct {
	DepartmentName   string   `json:"department_name"`
	ManageAccount    string   `json:"manage_account"`
	RoleList         []string `json:"role_list"`
	ParentDepartment string   `json:"parent_department"`
}

func (a ApprovalDetail) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-query approval_detail=======================")
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
	tabApproval := approvals.ApprovalView{}
	var resp ApprovalDetailResp
	dsApproval, err := tabApproval.QueryByFilter(global.DBMysql, map[string]interface{}{"approval_id": a.ApprovalId})
	logger.Info("[INFO] " + fmt.Sprintf("dsApproval: %v", dsApproval))
	if err != nil {
		logger.Error(err.Error())
		return ico.Err(2220, "", err.Error())
	}
	if dsApproval.ApprovalId == 0 {
		return ico.Err(2214, "", errors.New("the query approval data is empty"))
	}

	resp.UserInfo.NewInfo.Account = ""
	resp.UserInfo.NewInfo.UserName = ""
	resp.UserInfo.NewInfo.Department = ""
	resp.UserInfo.NewInfo.ManageDepartment = nil
	resp.UserInfo.NewInfo.RoleName = ""
	resp.UserInfo.OldInfo.Account = ""
	resp.UserInfo.OldInfo.UserName = ""
	resp.UserInfo.OldInfo.Department = ""
	resp.UserInfo.OldInfo.ManageDepartment = nil
	resp.UserInfo.OldInfo.RoleName = ""

	resp.DepartmentInfo.NewInfo.DepartmentName = ""
	resp.DepartmentInfo.NewInfo.ManageAccount = ""
	resp.DepartmentInfo.NewInfo.RoleList = nil
	resp.DepartmentInfo.NewInfo.ParentDepartment = ""
	resp.DepartmentInfo.OldInfo.DepartmentName = ""
	resp.DepartmentInfo.OldInfo.ManageAccount = ""
	resp.DepartmentInfo.OldInfo.RoleList = nil
	resp.DepartmentInfo.OldInfo.ParentDepartment = ""

	switch dsApproval.ApprovalType {
	case 1:
		ua := public.UserAdd{}
		if err := json.Unmarshal(dsApproval.ApprovalMessage, &ua); err != nil {
			logger.Error(err)
		}
		resp.UserInfo.NewInfo.Account = ua.Account
		resp.UserInfo.NewInfo.UserName = ua.UserName
		resp.UserInfo.NewInfo.Department = a.getDeptName(deptId)
		resp.UserInfo.NewInfo.ManageDepartment = a.getDeptNames(ua.ManageDeptId)
		resp.UserInfo.NewInfo.RoleName = a.getRoleName(ua.RoleId)
	case 2:
		uu := public.UserUpdate{}
		if err := json.Unmarshal(dsApproval.ApprovalMessage, &uu); err != nil {
			logger.Error(err)
		}
		resp.UserInfo.NewInfo.UserName = uu.UserNameNew
		resp.UserInfo.NewInfo.ManageDepartment = a.getDeptNames(uu.ManageDeptIdNew)
		resp.UserInfo.NewInfo.RoleName = a.getRoleName(uu.RoleIdNew)
		resp.UserInfo.OldInfo.UserName = uu.UserName
		resp.UserInfo.OldInfo.ManageDepartment = a.getDeptNames(uu.ManageDeptId)
		resp.UserInfo.OldInfo.RoleName = a.getRoleName(uu.RoleId)
	case 3:
		da := public.DeptAdd{}
		if err := json.Unmarshal(dsApproval.ApprovalMessage, &da); err != nil {
			logger.Error(err)
		}
		resp.DepartmentInfo.NewInfo.DepartmentName = da.DeptName
		resp.DepartmentInfo.NewInfo.ManageAccount = a.getAccount(da.ChargerId)
		resp.DepartmentInfo.NewInfo.RoleList = a.getRoleNames(da.RoleId)
		resp.DepartmentInfo.NewInfo.ParentDepartment = a.getDeptName(da.ParentId)
	case 4:
		du := public.DeptUpdate{}
		if err := json.Unmarshal(dsApproval.ApprovalMessage, &du); err != nil {
			logger.Error(err)
		}
		resp.DepartmentInfo.NewInfo.DepartmentName = du.DeptNameNew
		resp.DepartmentInfo.NewInfo.ManageAccount = a.getAccount(du.ChargerIdNew)
		resp.DepartmentInfo.NewInfo.RoleList = a.getRoleNames(du.RoleIdNew)
		resp.DepartmentInfo.OldInfo.DepartmentName = du.DeptName
		resp.DepartmentInfo.OldInfo.ManageAccount = a.getAccount(du.ChargerId)
		resp.DepartmentInfo.OldInfo.RoleList = a.getRoleNames(du.RoleId)
	}
	resp.ApplicantInfo.Account = dsApproval.Account
	resp.ApplicantInfo.UserName = dsApproval.UserName
	resp.ApplicantInfo.Department = dsApproval.Department
	resp.ApplicantInfo.ManageDepartment = a.getManagerDept(dsApproval.ChargerId)
	resp.ApplicantInfo.RoleName = dsApproval.RoleName
	resp.SubmitTime = dsApproval.CreateTime

	resp.History = append(resp.History, History{
		Time:           dsApproval.CreateTime,
		UserName:       dsApproval.UserName,
		ApprovalStatus: dsApproval.ApprovalStatus,
		Message:        dsApproval.SubmitMsg,
	})

	if dsApproval.ApprovalStatus != 1 {
		resp.History = append(resp.History, History{
			Time:           dsApproval.ApprovalTime,
			UserName:       "root",
			ApprovalStatus: dsApproval.ApprovalStatus,
			Message:        dsApproval.ApprovalMsg,
		})
	}
	resp.ApprovalType = 1

	return ico.Succ(resp)
}

func (a ApprovalDetail) getManagerDept(chargerId int) (deptNames []string) {
	tabDept := departments.DepartmentTab{}
	dsDept, err := tabDept.QueryListByFilter(global.DBMysql, map[string]interface{}{"charger_id": chargerId})
	logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	if len(dsDept) == 0 {
		return nil
	}
	for _, v := range dsDept {
		deptNames = append(deptNames, v.DeptName)
	}

	return
}

func (a ApprovalDetail) getDeptName(deptId int) (deptName string) {
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

func (a ApprovalDetail) getDeptNames(deptIds []int) (deptNames []string) {
	tabDept := departments.DepartmentTab{}
	for _, id := range deptIds {
		dsDept, err := tabDept.QueryListByFilter(global.DBMysql, map[string]interface{}{"id": id})
		logger.Info("[INFO] " + fmt.Sprintf("dsDept: %v", dsDept))
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
		if len(dsDept) == 0 {
			return nil
		}
		for _, v := range dsDept {
			deptNames = append(deptNames, v.DeptName)
		}
	}

	return
}

func (a ApprovalDetail) getRoleName(roleId int) (roleNames string) {
	tabRole := roles.RoleTable{}
	dsRole, err := tabRole.QueryByFilter(global.DBMysql, map[string]interface{}{"id": roleId})
	logger.Info("[INFO] " + fmt.Sprintf("dsRole: %v", dsRole))
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	if dsRole.RoleId == 0 {
		return ""
	}
	roleNames = dsRole.RoleName

	return
}

func (a ApprovalDetail) getRoleNames(roleIds []int) (roleNames []string) {
	tabRole := roles.RoleTable{}
	for _, id := range roleIds {
		dsRole, err := tabRole.QueryByFilter(global.DBMysql, map[string]interface{}{"id": id})
		logger.Info("[INFO] " + fmt.Sprintf("dsRole: %v", dsRole))
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
		if dsRole.RoleId == 0 {
			return nil
		}
		roleNames = append(roleNames, dsRole.RoleName)
	}

	return
}

func (a ApprovalDetail) getAccount(userId int) (account string) {
	tabUser := users.UserTable{}
	dsUser, err := tabUser.QueryByFilter(global.DBMysql, map[string]interface{}{"id": userId})
	logger.Info("[INFO] " + fmt.Sprintf("dsUser: %v", dsUser))
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	if dsUser.UserId == 0 {
		return ""
	}
	account = dsUser.Account

	return
}

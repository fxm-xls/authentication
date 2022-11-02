package departments

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DeptQuery struct {
	Type         int    `json:"type"` // 是否分页（0：分页默认；1：不分页）
	PageSize     int    `json:"page_size"`
	PageIndex    int    `json:"page_index"`
	DepartmentId int    `json:"department_id"` // 部门id
	Department   string `json:"department"`    // 部门名称 模糊查找
	Managed      int    `json:"managed"`       // 是否被管理（0：全部默认；1：已有管理；2：没有管理）
}

type DeptQueryResp struct {
	Total          int64      `json:"total"`
	DepartmentList []DeptInfo `json:"department_list"`
}

type DeptInfo struct {
	DepartmentId         int        `json:"department_id"`
	DepartmentName       string     `json:"department_name"`
	CreateTime           int64      `json:"create_time"`
	DepartmentManager    string     `json:"department_manager"`
	DepartmentManagerId  int        `json:"department_manager_id"`
	RoleList             []RoleInfo `json:"role_list"`
	ParentDepartmentId   int        `json:"parent_department_id"`
	ParentDepartmentName string     `json:"parent_department_name"`
	DepartmentLevel      int        `json:"department_level"`
}

type RoleInfo struct {
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (d DeptQuery) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-query department=======================")
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

	tabDept := departments.DeptUserView{}
	var dsDept []departments.DeptUserView
	var resp DeptQueryResp
	var err error
	var cdn string
	data := make(map[string]interface{})
	if d.DepartmentId > 0 {
		data = map[string]interface{}{"department_id": d.DepartmentId}
	}

	switch d.Managed {
	case 0: // 全部默认
		cdn = ""
	case 1: // 已有管理
		cdn = fmt.Sprintf("department_manager_id = %d", userId)
	case 2: // 没有管理
		cdn = fmt.Sprintf("department_manager_id != %d OR department_manager_id is NULL OR department_manager_id = 0", userId)
	}

	if d.Type == 0 { // paging
		if d.PageIndex == 0 || d.PageSize == 0 {
			d.PageIndex = 1
			d.PageSize = 10
		}
		pageMsg := utils.Pagination{PageIndex: d.PageIndex, PageSize: d.PageSize}
		resp.Total, dsDept, err = tabDept.QueryMulti(global.DBMysql, pageMsg, data, d.Department, cdn)
		logger.Info("[INFO] " + fmt.Sprintf("[paging] count: %d, dsDept: %v", resp.Total, dsDept))
	} else { // no paging
		dsDept, err = tabDept.QueryListByFilter(global.DBMysql, data, d.Department, cdn)
		logger.Info("[INFO] " + fmt.Sprintf("[no paging] dsDept: %v", dsDept))
	}
	if err != nil {
		logger.Error(err.Error())
		return ico.Err(2202, "", err.Error())
	}
	for _, v := range dsDept {
		resp.DepartmentList = append(resp.DepartmentList, DeptInfo{
			DepartmentId:         v.DepartmentId,
			DepartmentName:       v.DepartmentName,
			CreateTime:           v.CreateTime,
			DepartmentManager:    v.DepartmentManager,
			DepartmentManagerId:  v.DepartmentManagerID,
			RoleList:             d.getRoles(v.DepartmentId),
			ParentDepartmentId:   v.ParentDepartmentId,
			ParentDepartmentName: v.ParentDepartmentName,
			DepartmentLevel:      v.DepartmentLevel,
		})
	}

	return ico.Succ(resp)
}

func (d DeptQuery) getRoles(deptId int) (roleInfo []RoleInfo) {
	tabDeptRole := departments.DeptRoleView{}
	dsDeptRole, err := tabDeptRole.QueryListByFilter(global.DBMysql, map[string]interface{}{"dept_id": deptId})
	logger.Info("[INFO] " + fmt.Sprintf("dsDeptRole: %v", dsDeptRole))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	if len(dsDeptRole) == 0 {
		return nil
	}
	for _, v := range dsDeptRole {
		roleInfo = append(roleInfo, RoleInfo{
			RoleId:   v.RoleId,
			RoleName: v.RoleName,
		})
	}

	return
}

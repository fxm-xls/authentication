package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/roles"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

type RoleQuery struct {
	Type         int    `json:"type"`
	PageIndex    int    `json:"page_index"`
	PageSize     int    `json:"page_size"`
	RoleName     string `json:"role_name"`
	DepartmentId int    `json:"department_id"`
}

type RoleQueryRes struct {
	Total    int64      `json:"total"`
	RoleList []RoleInfo `json:"role_list"`
}

type RoleInfo struct {
	RoleId     int    `json:"role_id"`
	RoleName   string `json:"role_name"`
	RoleDesc   string `json:"role_desc"`
	UserNum    int    `json:"user_num"`
	UserNumEd  int    `json:"user_num_ed"`
	RoleNum    int    `json:"role_num"`
	RoleNumEd  int    `json:"role_num_ed"`
	CreateTime int64  `json:"create_time"`
}

func (This RoleQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("角色查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.该用户是否是部门管理员
	if code, err := public.JudgeManager(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取单个服务
	res, code, err := This.GetRole(userId)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ(res)
}

func (This RoleQuery) GetRole(userId int) (res RoleQueryRes, code int, err error) {
	// 1.获取角色
	roleT := roles.RoleTable{}
	whereMap := map[string]interface{}{"service_id": global.DefaultServiceId}
	// 分页
	var roleListQuery []roles.RoleTable
	if This.Type == 0 {
		if This.PageIndex == 0 || This.PageSize == 0 {
			return res, 2099, errors.New("参数异常")
		}
		pageMsg := utils.Pagination{PageIndex: This.PageIndex, PageSize: This.PageSize}
		count, roleList, err := roleT.QueryMulti(global.DBMysql, pageMsg, whereMap, This.RoleName)
		if err != nil {
			logger.Error("角色查询失败")
			return res, 2112, errors.New("")
		}
		roleListQuery = roleList
		res.Total = count
	} else {
		// 不分页
		roleList, err := roleT.QueryListByFilter(global.DBMysql, whereMap)
		if err != nil {
			logger.Error("角色查询失败")
			return res, 2112, errors.New("")
		}
		roleListQuery = roleList
		res.Total = int64(len(roleList))
	}
	// 2.部门过滤
	roleIds := []int{}
	if This.DepartmentId != 0 {
		deptRoleT := departments.DeptRoleTab{}
		departmentRoleList, err := deptRoleT.QueryListByFilter(global.DBMysql, map[string]interface{}{"dept_id": This.DepartmentId})
		if err != nil {
			logger.Error("角色查询失败")
			return res, 2112, errors.New("该部门无可用角色")
		}
		for _, departmentRole := range departmentRoleList {
			roleIds = append(roleIds, departmentRole.RoleId)
		}
		res.Total = int64(len(roleIds))
	}
	// 3.删去其他信息
	for _, role := range roleListQuery {
		if This.DepartmentId != 0 && !utils.IsContainsInt(roleIds, role.RoleId) {
			continue
		}
		res.RoleList = append(res.RoleList, RoleInfo{
			RoleId:     role.RoleId,
			RoleName:   role.RoleName,
			RoleDesc:   role.RoleDesc,
			UserNum:    role.UserNum,
			RoleNum:    role.RoleNum,
			CreateTime: role.CreateTime,
		})
	}
	// 4.计算绑定值
	RoleListEd, code, err := This.GetBindNum(res.RoleList)
	if err != nil {
		return res, code, err
	}
	res.RoleList = RoleListEd
	return
}

func (This RoleQuery) GetBindNum(roleList []RoleInfo) (roleListEd []RoleInfo, code int, err error) {
	// 1.查询所有user-role
	userRoleT := users.UserRoleTable{}
	userRoleList, err := userRoleT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error("用户与角色联系查询失败")
		return roleListEd, 2142, errors.New("")
	}
	// 2.查询所有role-sub_role
	roleSubT := roles.RoleSubTable{}
	roleSubList, err := roleSubT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error("角色与大角色联系查询失败")
		return roleListEd, 2118, errors.New("")
	}
	// 3.计算绑定值
	ruleListJson, _ := json.Marshal(roleList)
	_ = json.Unmarshal(ruleListJson, &roleListEd)
	for i, role := range roleList {
		// 3.1 user-num
		userNum := 0
		for _, userRole := range userRoleList {
			if userRole.RoleId == role.RoleId {
				userNum++
			}
		}
		roleListEd[i].UserNumEd = userNum
		// 3.1 role-num
		roleNum := 0
		for _, roleSub := range roleSubList {
			if roleSub.SubRoleId == role.RoleId {
				roleNum++
			}
		}
		roleListEd[i].RoleNumEd = roleNum
	}
	return
}

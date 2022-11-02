package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/departments"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserQuery struct {
	Type             int    `json:"type"`
	PageIndex        int    `json:"page_index"`
	PageSize         int    `json:"page_size"`
	RoleId           int    `json:"role_id"`
	DepartmentId     int    `json:"department_id"`
	Department       string `json:"department"`
	UserName         string `json:"user_name"`
	DepartmentIdList []int
}

type UserQueryRes struct {
	Total    int64      `json:"total"`
	UserList []UserInfo `json:"user_list"`
}

type UserInfo struct {
	UserId               int      `json:"user_id"`
	Account              string   `json:"account"`
	UserName             string   `json:"user_name"`
	CreateTime           int64    `json:"create_time"`
	RoleId               int      `json:"role_id"`
	RoleName             string   `json:"role_name"`
	DepartmentId         int      `json:"department_id"`
	DepartmentName       string   `json:"department_name"`
	ManageDepartmentId   []int    `json:"manage_department_id"`
	ManageDepartmentName []string `json:"manage_department_name"`
}

func (This UserQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("用户查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.该用户是否是是部门管理员
	if _, err := public.JudgeManager(userId); err != nil {
		res, code, err := This.GetUserSelf(userId)
		if err != nil {
			return ico.Err(code, err.Error())
		}
		return ico.Succ(res)
	}
	// 2.获取部门id列表
	code, err := This.GetManageDepartmentId(userId)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.获取用户
	res, code, err := This.GetUsers()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ(res)
}

func (This UserQuery) GetUserSelf(userId int) (res UserQueryRes, code int, err error) {
	userV := users.UserView{}
	user, err := userV.QueryByFilter(global.DBMysql, map[string]interface{}{"user_id": userId})
	if err != nil {
		logger.Error("用户不存在")
		return res, 2144, errors.New("")
	}
	res.Total = 1
	res.UserList = append(res.UserList, UserInfo{
		UserId: user.UserId, UserName: user.UserName, Account: user.Account, CreateTime: user.CreateTime,
		RoleId: user.RoleId, RoleName: user.RoleName, DepartmentName: user.Department, DepartmentId: user.DepartmentId,
		ManageDepartmentId: []int{}, ManageDepartmentName: []string{},
	})
	return
}

func (This *UserQuery) GetManageDepartmentId(userId int) (code int, err error) {
	departmentT := departments.DepartmentTab{}
	whereMap := map[string]interface{}{"charger_id": userId}
	departmentList, err := departmentT.QueryListByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Error(err)
		return 2148, errors.New("")
	}
	for _, department := range departmentList {
		This.DepartmentIdList = append(This.DepartmentIdList, department.Id)
	}
	return
}

func (This UserQuery) GetUsers() (res UserQueryRes, code int, err error) {
	userV := users.UserView{}
	var userViewList []users.UserView
	// 0.条件过滤
	whereMap := map[string]interface{}{}
	if This.RoleId != 0 {
		whereMap["role_id"] = This.RoleId
	}
	if This.DepartmentId != 0 {
		whereMap["dept_id"] = This.DepartmentId
	}
	// 1.分页
	if This.Type == 0 {
		if This.PageIndex == 0 || This.PageSize == 0 {
			logger.Error("参数异常")
			return res, 2099, errors.New("分页异常")
		}
		pageMsg := utils.Pagination{PageIndex: This.PageIndex, PageSize: This.PageSize}
		defaultInt := []int{0, 2}
		res.Total, userViewList, err = userV.QueryMulti(global.DBMysql, pageMsg, defaultInt, This.UserName, This.Department, whereMap)
		if err != nil {
			logger.Error("用户查询失败")
			return res, 2135, errors.New("")
		}
	} else {
		// 2.不分页
		userViewList, err = userV.QueryListByFilter(global.DBMysql, map[string]interface{}{})
		if err != nil {
			logger.Error("用户查询失败")
			return res, 2135, errors.New("")
		}
		res.Total = int64(len(userViewList))
	}
	// 3.查询部门信息
	departmentT := departments.DepartmentTab{}
	departmentList, err := departmentT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error(err)
		return res, 2148, errors.New("")
	}
	// 4.整合信息
	for _, userView := range userViewList {
		userInfo := UserInfo{
			UserId: userView.UserId, UserName: userView.UserName, Account: userView.Account,
			CreateTime: userView.CreateTime, RoleId: userView.RoleId, RoleName: userView.RoleName,
			DepartmentId: userView.DepartmentId, DepartmentName: userView.Department,
			ManageDepartmentId: []int{}, ManageDepartmentName: []string{},
		}
		for _, department := range departmentList {
			if department.ChargerId == userView.UserId {
				userInfo.ManageDepartmentId = append(userInfo.ManageDepartmentId, department.Id)
				userInfo.ManageDepartmentName = append(userInfo.ManageDepartmentName, department.DeptName)
			}
		}
		res.UserList = append(res.UserList, userInfo)
	}
	return
}

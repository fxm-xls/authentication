package approvals

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/approvals"
	"github.com/gin-gonic/gin"
)

type ApprovalRevoke struct {
	List []RevokeList `json:"list"`
}

type RevokeList struct {
	ApprovalId int `json:"approval_id"`
}

func (a ApprovalRevoke) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-query approval_revoke=======================")
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
	var tx = global.DBMysql.Begin()
	for _, v := range a.List {
		tabApproval := approvals.ApprovalTab{Id: v.ApprovalId}
		data := map[string]interface{}{
			"status": 4,
		}
		if err := tabApproval.Update(global.DBMysql, v.ApprovalId, data); err != nil {
			logger.Error(err.Error())
			tx.Rollback()
			return ico.Err(2204, "", err.Error())
		}
	}
	tx.Commit()
	return ico.Succ("申请撤回成功")
}

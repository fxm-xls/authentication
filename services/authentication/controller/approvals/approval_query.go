package approvals

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/controller/public"
	"bigrule/services/authentication/model/approvals"
	"bigrule/services/authentication/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApprovalQuery struct {
	Type           int `json:"type"` // 是否分页（0：分页默认；1：不分页）
	PageSize       int `json:"page_size"`
	PageIndex      int `json:"page_index"`
	ApprovalStatus int `json:"approval_status"` // 申请状态（0：全部默认，1：审批中，2：已通过，3：已驳回，4：已撤回）
}

type ApprovalQueryResp struct {
	Total        int64          `json:"total"`
	ApprovalList []ApprovalInfo `json:"approval_list"`
}

type ApprovalInfo struct {
	ApprovalId      int    `json:"approval_id"`
	Account         string `json:"account"`
	Department      string `json:"department"`
	ApprovalType    int    `json:"approval_type"`    // 申请类型：1：新增账号，2：修改账号，3：新增部门，4：修改部门
	ApprovalMessage string `json:"approval_message"` // 申请账号/部门
	CreateTime      int64  `json:"create_time"`      // 十位时间戳
	ApprovalStatus  int    `json:"approval_status"`  // 申请状态：1：审批中，2：已审批，3：已驳回，4：已撤回
}

func (a ApprovalQuery) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("[INFO] =======================UserCenter: v3-query approval=======================")
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
	var dsApproval []approvals.ApprovalView
	var resp ApprovalQueryResp
	var err error
	data := make(map[string]interface{})
	if a.ApprovalStatus > 0 {
		data = map[string]interface{}{"approval_status": a.ApprovalStatus}
	}

	if a.Type == 0 { // paging
		if a.PageIndex == 0 || a.PageSize == 0 {
			a.PageIndex = 1
			a.PageSize = 10
		}
		pageMsg := utils.Pagination{PageIndex: a.PageIndex, PageSize: a.PageSize}
		resp.Total, dsApproval, err = tabApproval.QueryMulti(global.DBMysql, pageMsg, data)
		logger.Info("[INFO] " + fmt.Sprintf("[paging] count: %d, dsDept: %v", resp.Total, dsApproval))
	} else { // no paging
		dsApproval, err = tabApproval.QueryListByFilter(global.DBMysql, data)
		resp.Total = int64(len(dsApproval))
		logger.Info("[INFO] " + fmt.Sprintf("[no paging] dsDept: %v", dsApproval))
	}
	if err != nil {
		logger.Error(err.Error())
		return ico.Err(2220, "", err.Error())
	}

	var approvalMsg string
	for _, v := range dsApproval {
		switch v.ApprovalType {
		case 1: // 新增账号
			ua := public.UserAdd{}
			if err := json.Unmarshal(v.ApprovalMessage, &ua); err != nil {
				logger.Error(err)
			}
			approvalMsg = ua.UserName
		case 2: // 修改账号
			uu := public.UserUpdate{}
			if err := json.Unmarshal(v.ApprovalMessage, &uu); err != nil {
				logger.Error(err)
			}
			approvalMsg = uu.UserNameNew
		case 3: // 新增部门
			da := public.DeptAdd{}
			if err := json.Unmarshal(v.ApprovalMessage, &da); err != nil {
				logger.Error(err)
			}
			approvalMsg = da.DeptName
		case 4: // 修改部门
			du := public.DeptUpdate{}
			if err := json.Unmarshal(v.ApprovalMessage, &du); err != nil {
				logger.Error(err)
			}
			approvalMsg = du.DeptNameNew
		}

		resp.ApprovalList = append(resp.ApprovalList, ApprovalInfo{
			ApprovalId:      v.ApprovalId,
			Account:         v.Account,
			Department:      v.Department,
			ApprovalType:    v.ApprovalType,
			ApprovalMessage: approvalMsg,
			CreateTime:      v.CreateTime,
			ApprovalStatus:  v.ApprovalStatus,
		})
	}

	return ico.Succ(resp)
}

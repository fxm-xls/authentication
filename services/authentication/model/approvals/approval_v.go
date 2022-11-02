package approvals

import (
	"bigrule/services/authentication/utils"
	"gorm.io/gorm"
)

type ApprovalView struct {
	ApprovalId      int    `json:"approval_id"`
	Account         string `json:"account"`
	UserName        string `json:"user_name"`
	Department      string `json:"department"`
	ChargerId       int    `json:"charger_id"`
	RoleName        string `json:"role_name"`
	ApprovalType    int    `json:"approval_type"`
	ApprovalMessage JSON   `json:"approval_message"`
	CreateTime      int64  `json:"create_time"`
	ApprovalTime    int64  `json:"approval_time"`
	ApprovalStatus  int    `json:"approval_status"`
	SubmitMsg       string `json:"submit_msg"`
	ApprovalMsg     string `json:"approval_msg"`
}

func (ApprovalView) TableName() string {
	return "approval_v"
}

func (a *ApprovalView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (res ApprovalView, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).First(&res).Error
	return
}

func (a *ApprovalView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (res []ApprovalView, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).Order("create_time DESC").Find(&res).Error
	return
}

func (a *ApprovalView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, res []ApprovalView, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).Order("create_time DESC").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	return
}

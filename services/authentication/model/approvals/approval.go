package approvals

import (
	"bigrule/services/authentication/utils"
	"bytes"
	"database/sql/driver"
	"errors"
	"gorm.io/gorm"
)

type JSON []byte

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		_ = errors.New("invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

type ApprovalTab struct {
	Id           int    `json:"id"             gorm:"column:id;primary_key"`
	ApplicantId  int    `json:"applicant_id"   gorm:"column:applicant_id"`
	DeptId       int    `json:"dept_id"        gorm:"column:dept_id"`
	Type         int    `json:"type"           gorm:"column:type"`
	Content      JSON   `json:"content"        gorm:"column:content"`
	Status       int    `json:"status"         gorm:"column:status"`
	SubmitMsg    string `json:"submit_msg"     gorm:"column:submit_msg"`
	ApprovalMsg  string `json:"approval_msg"   gorm:"column:approval_msg"`
	CreateTime   int64  `json:"create_time"    gorm:"column:create_time"`
	ApprovalTime int64  `json:"approval_time"  gorm:"column:approval_time"`
}

func (ApprovalTab) TableName() string {
	return "approval_t"
}

func (a *ApprovalTab) Insert(tx *gorm.DB) error {
	err := tx.Table(a.TableName()).Create(a).Debug().Error
	if err == nil {
		err = tx.Table(a.TableName()).Last(a).Error
	}
	return err
}

func (a *ApprovalTab) Update(tx *gorm.DB, id int, data map[string]interface{}) error {
	err := tx.Table(a.TableName()).Debug().Where("id=?", id).Updates(data).Error
	return err
}

func (a *ApprovalTab) QueryByFilterPage(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, approval []ApprovalTab, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&approval).Error
	return
}

func (a *ApprovalTab) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (approval ApprovalTab, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).First(&approval).Error
	return
}

func (a *ApprovalTab) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (approval []ApprovalTab, err error) {
	err = tx.Table(a.TableName()).Debug().Where(data).Find(&approval).Error
	return
}

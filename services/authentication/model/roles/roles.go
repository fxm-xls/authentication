package roles

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type RoleTable struct {
	RoleId     int    `json:"role_id,omitempty"     gorm:"column:id;primary_key"`
	RoleName   string `json:"role_name"             gorm:"column:name"`
	RoleDesc   string `json:"role_desc"             gorm:"column:desc"`
	ServiceId  int    `json:"service_id"            gorm:"column:service_id"`
	Manager    int    `json:"manager"               gorm:"column:manager"`
	UserNum    int    `json:"user_num"              gorm:"column:user_num"`
	RoleNum    int    `json:"role_num"              gorm:"column:role_num"`
	Status     int    `json:"-"                     gorm:"column:status"`
	CreateTime int64  `json:"create_time"           gorm:"column:create_time"`
}

func (RoleTable) TableName() string {
	return "role_t"
}

func (s *RoleTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *RoleTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleId).Delete(nil).Error
	return err
}

func (s *RoleTable) DeleteByRoleIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("id in ?", data).Delete(nil).Error
	return err
}

func (s *RoleTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleId).Updates(upInfo).Error
	return err
}

func (s *RoleTable) UpdateByStruct(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Model(s).Updates(s).Error
	return err
}

func (s *RoleTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (role RoleTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&role).Error
	return
}

func (s *RoleTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (role []RoleTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&role).Error
	return
}

func (s *RoleTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, roleName string) (count int64, res []RoleTable, err error) {
	if roleName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", roleName)).Order("id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Order("id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

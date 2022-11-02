package roles

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type RoleSubTable struct {
	RoleSubId int `json:"id,omitempty"   gorm:"column:id;primary_key"`
	RoleId    int `json:"role_id"        gorm:"column:role_id"`
	SubRoleId int `json:"sub_role_id"    gorm:"column:sub_role_id"`
}

func (RoleSubTable) TableName() string {
	return "role_sub_t"
}

func (s *RoleSubTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *RoleSubTable) InsertMany(tx *gorm.DB, roleSubList []RoleSubTable) error {
	err := tx.Table(s.TableName()).Create(&roleSubList).Error
	return err
}

func (s *RoleSubTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleSubId).Delete(nil).Error
	return err
}

func (s *RoleSubTable) DeleteByRoleIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("role_id in ?", data).Delete(nil).Error
	return err
}

func (s *RoleSubTable) DeleteBySubRoleIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("sub_role_id in ?", data).Delete(nil).Error
	return err
}

func (s *RoleSubTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleSubId).Updates(upInfo).Error
	return err
}

func (s *RoleSubTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (roleSub RoleSubTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&roleSub).Error
	return
}

func (s *RoleSubTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (roleSub []RoleSubTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&roleSub).Error
	return
}

func (s *RoleSubTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, RoleSubName string) (count int64, res []RoleSubTable, err error) {
	if RoleSubName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", RoleSubName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

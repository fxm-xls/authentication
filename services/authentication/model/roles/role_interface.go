package roles

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type RoleInterfaceTable struct {
	RoleInterfaceId int `json:"role_interface_id,omitempty"   gorm:"column:id;primary_key"`
	RoleId          int `json:"role_id"                       gorm:"column:role_id"`
	InterfaceId     int `json:"interface_id"                  gorm:"column:interface_id"`
}

func (RoleInterfaceTable) TableName() string {
	return "role_interface_t"
}

func (s *RoleInterfaceTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *RoleInterfaceTable) InsertMany(tx *gorm.DB, roleInterfaceList []RoleInterfaceTable) error {
	err := tx.Table(s.TableName()).Create(&roleInterfaceList).Error
	return err
}

func (s *RoleInterfaceTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleInterfaceId).Delete(nil).Error
	return err
}

func (s *RoleInterfaceTable) DeleteByRoleIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("role_id in ?", data).Delete(nil).Error
	return err
}

func (s *RoleInterfaceTable) DeleteByInterfaceIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("interface_id in ?", data).Delete(nil).Error
	return err
}

func (s *RoleInterfaceTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.RoleInterfaceId).Updates(upInfo).Error
	return err
}

func (s *RoleInterfaceTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (roleInterface RoleInterfaceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&roleInterface).Error
	return
}

func (s *RoleInterfaceTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (roleInterface []RoleInterfaceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&roleInterface).Error
	return
}

func (s *RoleInterfaceTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, RoleInterfaceName string) (count int64, res []RoleInterfaceTable, err error) {
	if RoleInterfaceName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", RoleInterfaceName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

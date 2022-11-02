package menus

import (
	"bigrule/services/authentication/utils"
	"gorm.io/gorm"
)

type InterfaceTable struct {
	InterfaceId   int    `json:"interface_id,omitempty"     gorm:"column:interface_id;primary_key"`
	InterfaceName string `json:"interface_name"             gorm:"column:interface_name"`
	GeneralAccess int    `json:"general_access"             gorm:"column:general_access"`
	Path          string `json:"path"                       gorm:"column:path"`
	Method        string `json:"method"                     gorm:"column:method"`
}

func (InterfaceTable) TableName() string {
	return "interface_t"
}

func (m *InterfaceTable) Insert(tx *gorm.DB) error {
	err := tx.Table(m.TableName()).Create(m).Error
	if err == nil {
		err = tx.Table(m.TableName()).Last(m).Error
	}
	return err
}

func (m *InterfaceTable) Delete(tx *gorm.DB) error {
	err := tx.Table(m.TableName()).Where("interface_id=?", m.InterfaceId).Delete(nil).Error
	return err
}

func (m *InterfaceTable) Update(tx *gorm.DB, gpInfo map[string]interface{}) error {
	err := tx.Table(m.TableName()).Where("interface_id=?", m.InterfaceId).Updates(gpInfo).Error
	return err
}

func (m *InterfaceTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (interfaceList []InterfaceTable, err error) {
	err = tx.Table(m.TableName()).Where(data).Find(&interfaceList).Error
	return interfaceList, err
}

func (m *InterfaceTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, res []InterfaceTable, err error) {
	err = tx.Table(m.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	return
}

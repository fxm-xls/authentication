package interfaces

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type InterfaceTable struct {
	InterfaceId   int    `json:"interface_id,omitempty"  gorm:"column:id;primary_key"`
	InterfaceName string `json:"interface_name"          gorm:"column:name"`
	InterfaceDesc string `json:"interface_desc"          gorm:"column:desc"`
	Path          string `json:"path"                    gorm:"column:path"`
	Method        string `json:"method"                  gorm:"column:method"`
	ServiceId     int    `json:"service_id"              gorm:"column:service_id"`
	Status        int    `json:"-"                       gorm:"column:status"`
}

func (InterfaceTable) TableName() string {
	return "interface_t"
}

func (s *InterfaceTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *InterfaceTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.InterfaceId).Delete(nil).Error
	return err
}

func (s *InterfaceTable) DeleteByServiceId(tx *gorm.DB, serviceId int) error {
	err := tx.Table(s.TableName()).Where("service_id = ?", serviceId).Delete(nil).Error
	return err
}

func (s *InterfaceTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.InterfaceId).Updates(upInfo).Error
	return err
}

func (s *InterfaceTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (interfaceInfo InterfaceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&interfaceInfo).Error
	return
}

func (s *InterfaceTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (interfaceInfo []InterfaceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&interfaceInfo).Error
	return
}

func (s *InterfaceTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, InterfaceName string) (count int64, res []InterfaceTable, err error) {
	if InterfaceName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", InterfaceName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

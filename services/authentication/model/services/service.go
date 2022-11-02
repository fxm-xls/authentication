package services

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type ServiceTable struct {
	ServiceId   int    `json:"service_id,omitempty"     gorm:"column:id;primary_key"`
	ServiceName string `json:"service_name"             gorm:"column:name"`
	ServiceDesc string `json:"service_desc"             gorm:"column:desc"`
	IndexUrl    string `json:"index_url"                gorm:"column:index_url"`
	Status      int    `json:"status"                   gorm:"column:status"`
}

func (ServiceTable) TableName() string {
	return "service_t"
}

func (s *ServiceTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *ServiceTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.ServiceId).Delete(nil).Error
	return err
}

func (s *ServiceTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.ServiceId).Updates(upInfo).Error
	return err
}

func (s *ServiceTable) UpdateByStruct(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Model(s).Updates(s).Error
	return err
}

func (s *ServiceTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (service ServiceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&service).Error
	return
}

func (s *ServiceTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (service []ServiceTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&service).Error
	return
}

func (s *ServiceTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, ServiceName string) (count int64, res []ServiceTable, err error) {
	if ServiceName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", ServiceName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

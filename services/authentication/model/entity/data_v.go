package entity

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type DataView struct {
	UserId      int    `json:"user_id"            gorm:"column:user_id"`
	Account     string `json:"account"            gorm:"column:account"`
	UserName    string `json:"user_name"          gorm:"column:user_name"`
	EntityId    int    `json:"entity_id"          gorm:"column:entity_id"`
	DataId      int    `json:"data_id"            gorm:"column:data_id"`
	DataType    string `json:"data_type"          gorm:"column:data_type"`
	Operation   string `json:"operation"          gorm:"column:operation"`
	ServiceId   int    `json:"service_id"         gorm:"column:service_id"`
	ServiceName string `json:"service_name"       gorm:"column:service_name"`
}

func (DataView) TableName() string {
	return "data_v"
}

func (s *DataView) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *DataView) InsertMany(tx *gorm.DB, entityList []DataView) error {
	err := tx.Table(s.TableName()).Create(&entityList).Error
	return err
}

func (s *DataView) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.EntityId).Delete(nil).Error
	return err
}

func (s *DataView) DeleteByIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("id in ?", data).Delete(nil).Error
	return err
}

func (s *DataView) DeleteByServiceId(tx *gorm.DB, serviceId int) error {
	err := tx.Table(s.TableName()).Where("service_id = ?", serviceId).Delete(nil).Error
	return err
}

func (s *DataView) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.EntityId).Updates(upInfo).Error
	return err
}

func (s *DataView) UpdateByStruct(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Model(s).Updates(s).Error
	return err
}

func (s *DataView) QueryByFilter(tx *gorm.DB, entity map[string]interface{}) (entityInfo DataView, err error) {
	err = tx.Table(s.TableName()).Where(entity).First(&entityInfo).Error
	return
}

func (s *DataView) QueryListByFilter(tx *gorm.DB, entity map[string]interface{}) (entityInfo []DataView, err error) {
	err = tx.Table(s.TableName()).Where(entity).Find(&entityInfo).Error
	return
}

func (s *DataView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, entity map[string]interface{}, EntityName string) (count int64, res []DataView, err error) {
	if EntityName != "" {
		err = tx.Table(s.TableName()).Where(entity).Where(fmt.Sprintf("name like '%%%s%%'", EntityName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(entity).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

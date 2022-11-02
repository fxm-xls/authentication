package entity

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type EntityTable struct {
	EntityId  int    `json:"id,omitempty"       gorm:"column:id;primary_key"`
	DataId    int    `json:"data_id"            gorm:"column:data_id"`
	DataType  string `json:"data_type"          gorm:"column:data_type"`
	Operation string `json:"operation"          gorm:"column:operation"`
	DataName  string `json:"data_name"          gorm:"column:name"`
	DataDesc  string `json:"data_desc"          gorm:"column:desc"`
	ServiceId int    `json:"service_id"         gorm:"column:service_id"`
	Status    int    `json:"-"                  gorm:"column:status"`
}

func (EntityTable) TableName() string {
	return "entity_t"
}

func (s *EntityTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *EntityTable) InsertMany(tx *gorm.DB, entityList []EntityTable) error {
	err := tx.Table(s.TableName()).Create(&entityList).Error
	return err
}

func (s *EntityTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.EntityId).Delete(nil).Error
	return err
}

func (s *EntityTable) DeleteByIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("id in ?", data).Delete(nil).Error
	return err
}

func (s *EntityTable) DeleteByServiceId(tx *gorm.DB, serviceId int) error {
	err := tx.Table(s.TableName()).Where("service_id = ?", serviceId).Delete(nil).Error
	return err
}

func (s *EntityTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.EntityId).Updates(upInfo).Error
	return err
}

func (s *EntityTable) UpdateByStruct(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Model(s).Updates(s).Error
	return err
}

func (s *EntityTable) QueryByFilter(tx *gorm.DB, entity map[string]interface{}) (entityInfo EntityTable, err error) {
	err = tx.Table(s.TableName()).Where(entity).First(&entityInfo).Error
	return
}

func (s *EntityTable) QueryListByFilter(tx *gorm.DB, entity map[string]interface{}) (entityInfo []EntityTable, err error) {
	err = tx.Table(s.TableName()).Where(entity).Find(&entityInfo).Error
	return
}

func (s *EntityTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, entity map[string]interface{}, EntityName string) (count int64, res []EntityTable, err error) {
	if EntityName != "" {
		err = tx.Table(s.TableName()).Where(entity).Where(fmt.Sprintf("name like '%%%s%%'", EntityName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(entity).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

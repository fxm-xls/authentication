package users

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserEntityTable struct {
	UserEntityId int `json:"id,omitempty"    gorm:"column:id;primary_key"`
	UserId       int `json:"user_id"         gorm:"column:user_id"`
	EntityId     int `json:"entity_id"       gorm:"column:entity_id"`
}

func (UserEntityTable) TableName() string {
	return "user_entity_t"
}

func (s *UserEntityTable) Insert(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Create(s).Error
	if err == nil {
		err = tx.Table(s.TableName()).Last(s).Error
	}
	return err
}

func (s *UserEntityTable) InsertMany(tx *gorm.DB, userEntityList []UserEntityTable) error {
	err := tx.Table(s.TableName()).Create(&userEntityList).Error
	return err
}

func (s *UserEntityTable) Delete(tx *gorm.DB) error {
	err := tx.Table(s.TableName()).Where("id=?", s.UserEntityId).Delete(nil).Error
	return err
}

func (s *UserEntityTable) DeleteByEntityIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("entity_id in ?", data).Delete(nil).Error
	return err
}

func (s *UserEntityTable) DeleteByUserIds(tx *gorm.DB, data []int) error {
	err := tx.Table(s.TableName()).Where("user_id in ?", data).Delete(nil).Error
	return err
}

func (s *UserEntityTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(s.TableName()).Where("id=?", s.UserEntityId).Updates(upInfo).Error
	return err
}

func (s *UserEntityTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (userEntity UserEntityTable, err error) {
	err = tx.Table(s.TableName()).Where(data).First(&userEntity).Error
	return
}

func (s *UserEntityTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (userEntity []UserEntityTable, err error) {
	err = tx.Table(s.TableName()).Where(data).Find(&userEntity).Error
	return
}

func (s *UserEntityTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, UserEntityName string) (count int64, res []UserEntityTable, err error) {
	if UserEntityName != "" {
		err = tx.Table(s.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", UserEntityName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(s.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

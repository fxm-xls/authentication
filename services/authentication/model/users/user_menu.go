package users

import (
	"gorm.io/gorm"
)

type UserMenuTable struct {
	Id         int `json:"id,omitempty"     gorm:"column:id;primary_key"`
	UserId     int `json:"user_id"          gorm:"column:user_id"`
	MenuId     int `json:"menu_id"          gorm:"column:menu_id"`
	MenuType   int `json:"menu_type"        gorm:"column:menu_type"`
	MenuStatus int `json:"menu_status"      gorm:"column:menu_status"`
}

func (UserMenuTable) TableName() string {
	return "user_menu_t"
}

func (um *UserMenuTable) InsertMany(tx *gorm.DB, userMenuList []UserMenuTable) error {
	err := tx.Table(um.TableName()).Create(&userMenuList).Error
	return err
}

func (um *UserMenuTable) Delete(tx *gorm.DB) error {
	err := tx.Table(um.TableName()).Where("user_id=?", um.UserId).Delete(nil).Error
	return err
}

func (um *UserMenuTable) DeleteMany(tx *gorm.DB, MenuIdList []int) error {
	err := tx.Table(um.TableName()).Where("menu_id in ?", MenuIdList).Delete(nil).Error
	return err
}

func (um *UserMenuTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(um.TableName()).Where("user_id=?", um.UserId).Updates(upInfo).Error
	return err
}

func (um *UserMenuTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (userMenuList []UserMenuTable, err error) {
	err = tx.Table(um.TableName()).Where(data).Find(&userMenuList).Error
	return
}

func (um *UserMenuTable) QueryByRaw(tx *gorm.DB, sql string) (userMenuList []UserMenuTable, err error) {
	err = tx.Table(um.TableName()).Raw(sql).Find(&userMenuList).Error
	return
}

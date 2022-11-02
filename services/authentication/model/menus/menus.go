package menus

import (
	"bigrule/services/authentication/utils"
	"gorm.io/gorm"
)

type MenuTable struct {
	Id       int    `json:"id,omitempty"     gorm:"column:id;primary_key"`
	MenuType int    `json:"menu_type"        gorm:"column:menu_type"`
	DataId   int    `json:"data_id"          gorm:"column:data_id"`
	DataName string `json:"data_name"        gorm:"column:data_name"`
}

func (MenuTable) TableName() string {
	return "menu_t"
}

func (m *MenuTable) InsertMany(tx *gorm.DB, menuList []MenuTable) error {
	err := tx.Table(m.TableName()).Create(&menuList).Error
	return err
}

func (m *MenuTable) DeleteAll(tx *gorm.DB) error {
	err := tx.Table(m.TableName()).Where("id != 0").Delete(nil).Error
	return err
}

func (m *MenuTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (menuList []MenuTable, err error) {
	err = tx.Table(m.TableName()).Where(data).Find(&menuList).Error
	return menuList, err
}

func (m *MenuTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, res []MenuTable, err error) {
	err = tx.Table(m.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	return
}

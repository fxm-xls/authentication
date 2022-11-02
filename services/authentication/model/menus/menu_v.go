package menus

import "gorm.io/gorm"

type MenuView struct {
	MenuStatus int    `json:"menu_status"      gorm:"column:menu_status"`
	MenuType   int    `json:"menu_type"        gorm:"column:menu_type"`
	Path       string `json:"path"             gorm:"column:path"`
	Method     string `json:"method"           gorm:"column:method"`
}

func (MenuView) TableName() string {
	return "menu_v"
}

func (u *MenuView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (user MenuView, err error) {
	err = tx.Table(u.TableName()).Where(data).First(&user).Error
	return
}

func (u *MenuView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (user []MenuView, err error) {
	err = tx.Table(u.TableName()).Where(data).Find(&user).Error
	return
}

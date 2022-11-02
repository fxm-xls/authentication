package menus

import "gorm.io/gorm"

type MenuInterfaceTable struct {
	Id          int `json:"id,omitempty"     gorm:"column:id;primary_key"`
	MenuType    int `json:"menu_type"        gorm:"column:menu_type"`
	MenuStatus  int `json:"menu_status"      gorm:"column:menu_status"`
	InterfaceId int `json:"interface_id"     gorm:"column:interface_id"`
}

func (MenuInterfaceTable) TableName() string {
	return "priority_interface_t"
}

func (mi *MenuInterfaceTable) Insert(tx *gorm.DB) error {
	err := tx.Table(mi.TableName()).Create(mi).Error
	if err == nil {
		err = tx.Table(mi.TableName()).Last(mi).Error
	}
	return err
}

func (mi *MenuInterfaceTable) InsertMany(tx *gorm.DB, menuInterfaceList []MenuInterfaceTable) error {
	err := tx.Table(mi.TableName()).Create(&menuInterfaceList).Error
	return err
}

func (mi *MenuInterfaceTable) Delete(tx *gorm.DB) error {
	err := tx.Table(mi.TableName()).Where("id=?", mi.Id).Delete(nil).Error
	return err
}

func (mi *MenuInterfaceTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(mi.TableName()).Where("id=?", mi.Id).Updates(upInfo).Error
	return err
}

func (mi *MenuInterfaceTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (menuInterfaceList []MenuInterfaceTable, err error) {
	err = tx.Table(mi.TableName()).Where(data).Find(&menuInterfaceList).Error
	return
}

package departments

import "gorm.io/gorm"

type DeptRoleTab struct {
	Id     int `json:"id"        gorm:"column:id;primary_key"`
	DeptId int `json:"dept_id"   gorm:"column:dept_id"`
	RoleId int `json:"role_id"   gorm:"column:role_id"`
}

func (DeptRoleTab) TableName() string {
	return "dept_role_t"
}

func (d *DeptRoleTab) Insert(tx *gorm.DB) error {
	err := tx.Table(d.TableName()).Create(d).Error
	if err == nil {
		err = tx.Table(d.TableName()).Last(d).Error
	}
	return err
}

func (d *DeptRoleTab) Delete(tx *gorm.DB) error {
	err := tx.Table(d.TableName()).Debug().Where("dept_id=?", d.DeptId).Delete(nil).Error
	return err
}

func (d *DeptRoleTab) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (dept DeptRoleTab, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).First(&dept).Error
	return
}

func (d *DeptRoleTab) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (dept []DeptRoleTab, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).Find(&dept).Error
	return
}

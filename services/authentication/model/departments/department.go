package departments

import (
	"gorm.io/gorm"
)

type DepartmentTab struct {
	Id         int    `json:"id"          gorm:"column:id;primary_key"`
	DeptName   string `json:"dept_name"   gorm:"column:dept_name"`
	ParentId   int    `json:"parent_id"   gorm:"column:parent_id"`
	ChargerId  int    `json:"charger_id"  gorm:"column:charger_id"`
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`
	Desc       string `json:"desc"        gorm:"column:desc"`
	DeptLevel  int    `json:"dept_level"  gorm:"column:dept_level"`
}

func (DepartmentTab) TableName() string {
	return "department_t"
}

func (d *DepartmentTab) Insert(tx *gorm.DB) (int, error) {
	err := tx.Table(d.TableName()).Debug().Create(d).Error
	if err == nil {
		err = tx.Table(d.TableName()).Last(d).Error
	}
	return d.Id, err
}

func (d *DepartmentTab) Delete(tx *gorm.DB) error {
	err := tx.Table(d.TableName()).Debug().Where("id=?", d.Id).Delete(nil).Error
	return err
}

func (d *DepartmentTab) Update(tx *gorm.DB, data map[string]interface{}) error {
	err := tx.Table(d.TableName()).Debug().Where("id=?", d.Id).Updates(data).Error
	return err
}

func (d *DepartmentTab) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (dept DepartmentTab, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).First(&dept).Error
	return
}

func (d *DepartmentTab) QueryByFilter1(tx *gorm.DB, data map[string]interface{}, id int) (dept DepartmentTab, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).Where("id!=?", id).First(&dept).Error
	return
}

func (d *DepartmentTab) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (dept []DepartmentTab, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).Find(&dept).Error
	return
}

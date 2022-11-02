package departments

import (
	"bigrule/services/authentication/utils"
	"gorm.io/gorm"
)

type DeptRoleView struct {
	DeptId   int    `json:"dept_id"`
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (DeptRoleView) TableName() string {
	return "dept_role_v"
}

func (d *DeptRoleView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (deptRole DeptRoleView, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).First(&deptRole).Error
	return
}

func (d *DeptRoleView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (res []DeptRoleView, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).Find(&res).Error
	return
}

func (d *DeptRoleView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, res []DeptRoleView, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	return
}

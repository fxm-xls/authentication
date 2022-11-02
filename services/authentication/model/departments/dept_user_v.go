package departments

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type DeptUserView struct {
	DepartmentId         int    `json:"department_id"`
	DepartmentName       string `json:"department_name"`
	CreateTime           int64  `json:"create_time"`
	DepartmentManager    string `json:"department_manager"`
	DepartmentManagerID  int    `json:"department_manager_id"`
	RoleId               int    `json:"role_id"`
	RoleName             string `json:"role_name"`
	ParentDepartmentId   int    `json:"parent_department_id"`
	ParentDepartmentName string `json:"parent_department_name"`
	DepartmentLevel      int    `json:"department_level"`
}

func (DeptUserView) TableName() string {
	return "dept_user_v"
}

func (d *DeptUserView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (deptRole DeptUserView, err error) {
	err = tx.Table(d.TableName()).Debug().Where(data).First(&deptRole).Error
	return
}

func (d *DeptUserView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}, deptName, cdn string) (res []DeptUserView, err error) {
	if deptName != "" {
		err = tx.Table(d.TableName()).Debug().Where(data).Where(cdn).Where(fmt.Sprintf("department_name like '%%%s%%'", deptName)).Order("create_time DESC").Find(&res).Error
	} else {
		err = tx.Table(d.TableName()).Debug().Where(data).Where(cdn).Order("create_time DESC").Find(&res).Error
	}
	return
}

func (d *DeptUserView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, deptName, cdn string) (count int64, res []DeptUserView, err error) {
	if deptName != "" {
		err = tx.Table(d.TableName()).Debug().Where(data).Where(cdn).Where(fmt.Sprintf("department_name like '%%%s%%'", deptName)).Order("create_time DESC").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(d.TableName()).Debug().Where(data).Where(cdn).Order("create_time DESC").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

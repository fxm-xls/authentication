package users

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserView struct {
	UserId       int    `json:"user_id"       gorm:"column:user_id"`
	Account      string `json:"account"       gorm:"column:account"`
	UserName     string `json:"user_name"     gorm:"column:user_name"`
	UserDesc     string `json:"user_desc"     gorm:"column:user_desc"`
	CreateTime   int64  `json:"create_time"   gorm:"column:create_time"`
	Default      int    `json:"-"             gorm:"column:default"`
	RoleId       int    `json:"role_id"       gorm:"column:role_id"`
	RoleName     string `json:"role_name"     gorm:"column:role_name"`
	DepartmentId int    `json:"department_id" gorm:"column:dept_id"`
	Department   string `json:"department"    gorm:"column:dept_name"`
	Manager      int    `json:"manager"       gorm:"column:manager"`
	ServiceName  string `json:"service_name"  gorm:"column:service_name"`
}

func (UserView) TableName() string {
	return "user_v"
}

func (u *UserView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (user UserView, err error) {
	err = tx.Table(u.TableName()).Where(data).First(&user).Error
	return
}

func (u *UserView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (user []UserView, err error) {
	err = tx.Table(u.TableName()).Where(data).Find(&user).Error
	return
}

func (u *UserView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data []int, userName, deptName string, whereMap map[string]interface{}) (count int64, res []UserView, err error) {
	if userName != "" {
		if deptName != "" {
			err = tx.Table(u.TableName()).Where("`default` in ?", data).Where(fmt.Sprintf("user_name like '%%%s%%'", userName)).Where(fmt.Sprintf("dept_name like '%%%s%%'", deptName)).Where(whereMap).Order("user_id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
		} else {
			err = tx.Table(u.TableName()).Where("`default` in ?", data).Where(fmt.Sprintf("user_name like '%%%s%%'", userName)).Where(whereMap).Order("user_id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
		}
	} else {
		if deptName != "" {
			err = tx.Table(u.TableName()).Where("`default` in ?", data).Where(fmt.Sprintf("dept_name like '%%%s%%'", deptName)).Where(whereMap).Order("user_id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
		} else {
			err = tx.Table(u.TableName()).Where("`default` in ?", data).Where(whereMap).Order("user_id desc").Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
		}
	}
	return
}

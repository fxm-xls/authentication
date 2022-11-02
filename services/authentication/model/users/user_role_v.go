package users

import (
	"bigrule/services/authentication/utils"
	"gorm.io/gorm"
)

type UserRoleView struct {
	UserId      int    `json:"user_id"       gorm:"column:user_id"`
	Account     string `json:"account"       gorm:"column:account"`
	UserName    string `json:"user_name"     gorm:"column:user_name"`
	UserDesc    string `json:"user_desc"     gorm:"column:user_desc"`
	CreateTime  int64  `json:"create_time"   gorm:"column:create_time"`
	RoleId      int    `json:"role_id"       gorm:"column:role_id"`
	RoleName    string `json:"role_name"     gorm:"column:role_name"`
	Manager     int    `json:"manager"       gorm:"column:manager"`
	ServiceName string `json:"service_name"  gorm:"column:service_name"`
}

func (UserRoleView) TableName() string {
	return "user_role_v"
}

func (u *UserRoleView) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (user UserRoleView, err error) {
	err = tx.Table(u.TableName()).Where(data).First(&user).Error
	return
}

func (u *UserRoleView) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (user []UserRoleView, err error) {
	err = tx.Table(u.TableName()).Where(data).Find(&user).Error
	return
}

func (u *UserRoleView) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}) (count int64, res []UserRoleView, err error) {
	err = tx.Table(u.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	return
}

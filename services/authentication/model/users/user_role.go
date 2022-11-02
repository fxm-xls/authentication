package users

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserRoleTable struct {
	UserRoleId int `json:"id,omitempty"   gorm:"column:id;primary_key"`
	UserId     int `json:"user_id"        gorm:"column:user_id"`
	RoleId     int `json:"role_id"        gorm:"column:role_id"`
}

func (UserRoleTable) TableName() string {
	return "user_role_t"
}

func (ur *UserRoleTable) Insert(tx *gorm.DB) error {
	err := tx.Table(ur.TableName()).Create(ur).Error
	if err == nil {
		err = tx.Table(ur.TableName()).Last(ur).Error
	}
	return err
}

func (ur *UserRoleTable) InsertMany(tx *gorm.DB, userRoleList []UserRoleTable) error {
	err := tx.Table(ur.TableName()).Create(&userRoleList).Error
	return err
}

func (ur *UserRoleTable) Delete(tx *gorm.DB) error {
	err := tx.Table(ur.TableName()).Where("id=?", ur.UserRoleId).Delete(nil).Error
	return err
}

func (ur *UserRoleTable) DeleteByRoleIds(tx *gorm.DB, data []int) error {
	err := tx.Table(ur.TableName()).Where("role_id in ?", data).Delete(nil).Error
	return err
}

func (ur *UserRoleTable) DeleteByUserIds(tx *gorm.DB, data []int) error {
	err := tx.Table(ur.TableName()).Where("user_id in ?", data).Delete(nil).Error
	return err
}

func (ur *UserRoleTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(ur.TableName()).Where("id=?", ur.UserRoleId).Updates(upInfo).Error
	return err
}

func (ur *UserRoleTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (userRole UserRoleTable, err error) {
	err = tx.Table(ur.TableName()).Where(data).First(&userRole).Error
	return
}

func (ur *UserRoleTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (userRole []UserRoleTable, err error) {
	err = tx.Table(ur.TableName()).Where(data).Find(&userRole).Error
	return
}

func (ur *UserRoleTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, UserRoleName string) (count int64, res []UserRoleTable, err error) {
	if UserRoleName != "" {
		err = tx.Table(ur.TableName()).Where(data).Where(fmt.Sprintf("name like '%%%s%%'", UserRoleName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(ur.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

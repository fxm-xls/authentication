package users

import (
	"bigrule/services/authentication/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserTable struct {
	UserId       int    `json:"user_id,omitempty"     gorm:"column:id;primary_key"`
	Account      string `json:"account"               gorm:"column:account"`
	UserName     string `json:"user_name"             gorm:"column:user_name"`
	UserDesc     string `json:"desc"                  gorm:"column:desc"`
	Password     string `json:"password"              gorm:"column:password"`
	CreateTime   int64  `json:"create_time"           gorm:"column:create_time"`
	MobilePhone  string `json:"mobilephone"           gorm:"column:mobilephone"`
	Email        string `json:"email"                 gorm:"column:email"`
	Sex          string `json:"sex"                   gorm:"column:sex"`
	City         string `json:"city"                  gorm:"column:city"`
	Department   string `json:"department"            gorm:"column:department"`
	DepartmentId int    `json:"department_id"         gorm:"column:dept_id"`
	Default      int    `json:"-"                     gorm:"column:default"`
}

func (UserTable) TableName() string {
	return "user_t"
}

func (u *UserTable) Insert(tx *gorm.DB) error {
	err := tx.Table(u.TableName()).Create(u).Error
	if err == nil {
		err = tx.Table(u.TableName()).Last(u).Error
	}
	return err
}

func (u *UserTable) Delete(tx *gorm.DB) error {
	err := tx.Table(u.TableName()).Where("id=?", u.UserId).Delete(nil).Error
	return err
}

func (u *UserTable) DeleteByFilter(tx *gorm.DB, data map[string]interface{}) error {
	err := tx.Table(u.TableName()).Where(data).Delete(nil).Error
	return err
}

func (u *UserTable) Update(tx *gorm.DB, upInfo map[string]interface{}) error {
	err := tx.Table(u.TableName()).Where("id=?", u.UserId).Updates(upInfo).Error
	return err
}

func (u *UserTable) UpdateByStruct(tx *gorm.DB) error {
	err := tx.Table(u.TableName()).Model(u).Updates(u).Error
	return err
}

func (u *UserTable) QueryByFilter(tx *gorm.DB, data map[string]interface{}) (user UserTable, err error) {
	err = tx.Table(u.TableName()).Where(data).First(&user).Error
	return
}

func (u *UserTable) QueryListByFilter(tx *gorm.DB, data map[string]interface{}) (user []UserTable, err error) {
	err = tx.Table(u.TableName()).Where(data).Find(&user).Error
	return
}

func (u *UserTable) QueryMulti(tx *gorm.DB, pageMsg utils.Pagination, data map[string]interface{}, userName string) (count int64, res []UserTable, err error) {
	if userName != "" {
		err = tx.Table(u.TableName()).Where(data).Where(fmt.Sprintf("user_name like '%%%s%%'", userName)).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	} else {
		err = tx.Table(u.TableName()).Where(data).Count(&count).Offset(pageMsg.GetOffSet()).Limit(pageMsg.GetPageSize()).Find(&res).Error
	}
	return
}

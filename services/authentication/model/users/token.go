package users

import "gorm.io/gorm"

type TokenTable struct {
	UserId  int    `json:"user_id,omitempty"     gorm:"column:user_id;primary_key"`
	Token   string `json:"token"                 gorm:"column:token"`
	DueDate int64  `json:"due_date"              gorm:"column:due_date"`
}

func (TokenTable) TableName() string {
	return "token_t"
}

func (t *TokenTable) Insert(tx *gorm.DB) error {
	err := tx.Table(t.TableName()).Create(t).Error
	if err == nil {
		err = tx.Table(t.TableName()).Last(t).Error
	}
	return err
}

func (t *TokenTable) Delete(tx *gorm.DB) error {
	err := tx.Table(t.TableName()).Where("user_id=?", t.UserId).Delete(nil).Error
	return err
}

func (t *TokenTable) Update(tx *gorm.DB, tpInfo map[string]interface{}) error {
	err := tx.Table(t.TableName()).Where("user_id=?", t.UserId).Updates(tpInfo).Error
	return err
}

func (t *TokenTable) Query(tx *gorm.DB) (token TokenTable, err error) {
	err = tx.Table(t.TableName()).Where("user_id=?", t.UserId).First(&token).Error
	return token, err
}

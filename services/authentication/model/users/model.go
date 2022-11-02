package users

import (
	"errors"
	"gorm.io/gorm"
)

// SetToken 增加token
func SetToken(tx *gorm.DB, userId int, token string, expireTimes int64) (err error) {
	tokenT := TokenTable{
		UserId:  userId,
		Token:   token,
		DueDate: expireTimes,
	}
	// 查询是否存在token
	if _, err = tokenT.Query(tx); err != nil {
		// 增加token
		err = tokenT.Insert(tx)
		return nil
	}
	// 已存在，更新token
	updateMap := map[string]interface{}{
		"token":    token,
		"due_date": expireTimes,
	}
	err = tokenT.Update(tx, updateMap)
	return
}

// DelToken 删除token
func DelToken(tx *gorm.DB, userId int) (err error) {
	tokenT := TokenTable{UserId: userId}
	err = tokenT.Delete(tx)
	return
}

// GetUser 查询用户是否存在
func GetUser(tx *gorm.DB, userMap map[string]interface{}) (err error) {
	userT := UserTable{}
	userList, err := userT.QueryListByFilter(tx, userMap)
	if err != nil || len(userList) == 0 {
		return errors.New("该用户不存在")
	}
	return
}

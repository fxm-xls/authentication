package public

import (
	"bigrule/services/authentication/model"
	"errors"
)

// JudgeManager 权限判断 该用户是否是部门管理员
func JudgeManager(userId int) (code int, err error) {
	// 1.该用户是否为超级管理员角色
	if err = model.JudgeAuthManager(userId); err == nil {
		return
	}
	// 2.该用户是否是部门管理员
	if err = model.JudgeDepartmentManager(userId); err == nil {
		return
	}
	return 2171, errors.New("")
}

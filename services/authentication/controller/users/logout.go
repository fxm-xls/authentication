package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/users"
	"github.com/gin-gonic/gin"
)

type Logout struct{}

func (This Logout) DoHandle(c *gin.Context) *ico.Result {
	userId := c.GetInt("user_id")
	userName := c.GetString("user_name")
	logger.Infof("用户登出 user: id %d, name %s", userId, userName)
	// 删除对应token
	if err := users.DelToken(global.DBMysql, userId); err != nil {
		return ico.Err(2134, "")
	}
	return ico.Succ("登出成功")
}

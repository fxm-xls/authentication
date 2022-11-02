package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/config"
	"bigrule/services/authentication/middleware/jwt"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type Login struct {
	Account  string `json:"account"  binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (This Login) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "")
	}
	logger.Infof("用户登陆 user: name %s, password %s", This.Account, This.Password)
	// 1.验证用户名
	userT := &users.UserTable{}
	user, err := userT.QueryByFilter(global.DBMysql, map[string]interface{}{"account": This.Account})
	if err != nil {
		return ico.Err(2141, "")
	}
	// 2.验证密码是否相等
	isValid := utils.CheckPasswordHash(This.Password, user.Password)
	if !isValid {
		return ico.Err(2141, "")
	}
	// 3.生成token
	userMsg := make(map[string]interface{})
	userMsg["user_id"] = user.UserId
	userMsg["user_name"] = user.UserName
	userMsg["dept_id"] = user.DepartmentId
	expireTimes := time.Now().Add(config.CookieConfig.MysqlTime * time.Second).Unix()
	token, err := jwt.GenToken(userMsg, expireTimes)
	if err != nil {
		return ico.Err(2133, "")
	}
	// 4.mysql 存入token
	if err = users.SetToken(global.DBMysql, user.UserId, token, expireTimes); err != nil {
		return ico.Err(2133, "")
	}
	Res := LoginRes{
		UserId: user.UserId,
		Token:  token,
	}
	return ico.Succ(Res)
}

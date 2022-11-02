package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (UserRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/users", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/login", ico.Handler(Login{}))
		r.POST("/logout", ico.Handler(Logout{}))
		r.POST("/query", ico.Handler(UserQuery{}))
		r.POST("/token/query", ico.Handler(UserTokenQuery{}))
		//r.POST("/role/query", ico.Handler(UserRoleQuery{}))
		r.POST("/data/query", ico.Handler(UserDataQuery{}))
		r.POST("/add", ico.Handler(UserAdd{}))
		r.POST("/update", ico.Handler(UserUpdate{}))
		r.POST("/data/update", ico.Handler(UserDataUpdate{}))
		r.POST("/delete", ico.Handler(UserDelete{}))
	}
	rJf := router.Group(fmt.Sprintf("/v2/users")).Use(jwt.Jwt())
	{
		rJf.POST("/jf-service/login", ico.Handler(TokenLogin{}))
		rJf.POST("/jf-service/add", ico.JFHandler(UserJFAdd{}))
	}
}

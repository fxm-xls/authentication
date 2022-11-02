package auth

import (
	"bigrule/common/ico"
	"bigrule/services/authentication/controller/permission"
	"bigrule/services/authentication/controller/users"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (sr Router) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/v2/permission")).Use(jwt.Jwt())
	{
		r.POST("/interfaces/check", ico.Handler(permission.InterfaceCheck{}))
		r.POST("/data/get", ico.Handler(permission.DataGet{}))
	}
	r3 := router.Group("/v2/users").Use(jwt.Jwt())
	{
		r3.POST("/token/query", ico.Handler(users.UserTokenQuery{}))
		r3.POST("/login", ico.Handler(users.Login{}))
		r3.POST("/logout", ico.Handler(users.Logout{}))
		r3.POST("/query", ico.Handler(users.UserQuery{}))
		r3.POST("/update", ico.Handler(users.UserUpdate{}))
		r3.POST("/data/query", ico.Handler(users.UserDataQuery{}))
	}
	r2 := router.Group("/v2/user").Use(jwt.Jwt())
	{
		r2.POST("/permission/query", ico.Handler(PermissionQuery{}))
		r2.POST("/permission/update", ico.Handler(PermissionUpdate{}))
		r2.POST("/manager/query", ico.Handler(ManagerQuery{}))
	}
	r1 := router.Group("/v1/user").Use(jwt.Jwt())
	{
		r1.POST("/manager/query", ico.Handler(ManagerQuery{}))
	}
}

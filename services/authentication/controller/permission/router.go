package permission

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type PermissionRouter struct{}

func (PermissionRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/permission", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/interfaces/check", ico.Handler(InterfaceCheck{}))
		r.POST("/data/get", ico.Handler(DataGet{}))
	}
}

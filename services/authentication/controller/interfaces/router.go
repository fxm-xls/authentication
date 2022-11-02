package interfaces

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type InterfaceRouter struct{}

func (InterfaceRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/interfaces", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/query", ico.Handler(InterfaceQuery{}))
		r.POST("/import", ico.Handler(InterfaceImport{}))
	}
}

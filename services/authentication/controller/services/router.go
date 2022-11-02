package services

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ServiceRouter struct{}

func (ServiceRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/services", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/query", ico.Handler(ServiceQuery{}))
		r.POST("/register", ico.Handler(ServiceRegister{}))
		r.POST("/delete", ico.Handler(ServiceDelete{}))
		r.POST("/update", ico.Handler(ServiceUpdate{}))
	}
}

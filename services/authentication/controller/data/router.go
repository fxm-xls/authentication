package data

import (
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DataRouter struct{}

func (DataRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/v2/data")).Use(jwt.Jwt())
	{
		r.POST("/query", ico.Handler(DataQuery{}))
		r.POST("/import", ico.Handler(DataImport{}))
		r.POST("/user/query", ico.Handler(DataUserQuery{}))
	}
}

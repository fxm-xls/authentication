package roles

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type RoleRouter struct{}

func (RoleRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/roles", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/query", ico.Handler(RoleQuery{}))
		r.POST("/sub-role/query", ico.Handler(RoleSubQuery{}))
		r.POST("/add", ico.Handler(RoleAdd{}))
		r.POST("/delete", ico.Handler(RoleDelete{}))
		r.POST("/update", ico.Handler(RoleUpdate{}))
		//r.POST("/interface/update", ico.Handler(RoleInterfaceUpdate{}))
		r.POST("/sub-role/update", ico.Handler(RoleSubUpdate{}))
		r.POST("/import", ico.Handler(RoleImport{}))
	}
}

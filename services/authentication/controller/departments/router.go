package departments

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DepartmentRouter struct {
}

func (DepartmentRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/departments", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/add", ico.Handler(DeptAdd{}))
		r.POST("/update", ico.Handler(DeptUpdate{}))
		r.POST("/delete", ico.Handler(DeptDelete{}))
		r.POST("/query", ico.Handler(DeptQuery{}))
	}
}

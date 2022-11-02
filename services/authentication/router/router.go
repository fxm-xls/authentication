package router

import (
	"bigrule/common/router"
	"bigrule/services/authentication/config"
	"bigrule/services/authentication/controller/approvals"
	"bigrule/services/authentication/controller/auth-v1-v2"
	"bigrule/services/authentication/controller/data"
	"bigrule/services/authentication/controller/departments"
	"bigrule/services/authentication/controller/permission"
	"bigrule/services/authentication/controller/prometheus"
	"bigrule/services/authentication/controller/roles"
	"bigrule/services/authentication/controller/users"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
)

func RouterSetup() {
	router.RouterRegister(
		//services.ServiceRouter{},
		users.UserRouter{},
		roles.RoleRouter{},
		//interfaces.InterfaceRouter{},
		data.DataRouter{},
		permission.PermissionRouter{},
		auth.Router{},
		WebStaticRouter{},
		prometheus.Router{},
		departments.DepartmentRouter{},
		approvals.ApprovalRouter{},
	)
}

// WebStaticRouter 加载静态文件
type WebStaticRouter struct{}

func (ws WebStaticRouter) Router(router *gin.Engine) {
	//加载静态文件
	router.LoadHTMLGlob(config.WebStaticConfig.Path + "/index.html")
	_ = mime.AddExtensionType(".js", "application/javascript")
	router.StaticFS("/flowauth", gin.Dir(config.WebStaticConfig.Path, true))
	router.GET("/flowauth", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

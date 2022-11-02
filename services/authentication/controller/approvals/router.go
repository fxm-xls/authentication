package approvals

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/services/authentication/middleware/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApprovalRouter struct {
}

func (ApprovalRouter) Router(router *gin.Engine) {
	r := router.Group(fmt.Sprintf("/%s/approvals", global.Version)).Use(jwt.Jwt())
	{
		r.POST("/query", ico.Handler(ApprovalQuery{}))
		r.POST("/reject", ico.Handler(ApprovalReject{}))
		r.POST("/adopt", ico.Handler(ApprovalAdopt{}))
		r.POST("/detail/query", ico.Handler(ApprovalDetail{}))
		r.POST("/revoke", ico.Handler(ApprovalRevoke{}))
	}
}

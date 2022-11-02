package prometheus

import (
	"bigrule/services/authentication/middleware/jwt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct{}

func (This Router) Router(router *gin.Engine) {
	// 初始化 Prometheus
	newRouter := prometheus.NewRegistry()
	newRouter.MustRegister(jwt.RequestAllCount)
	newRouter.MustRegister(jwt.RequestTime)
	handler := promhttp.HandlerFor(newRouter, promhttp.HandlerOpts{})
	r := router.Group("/metrics").Use()
	{
		r.GET("", func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		})
	}
}

package jwt

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RequestAllCount 请求总次数
var RequestAllCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_method_counted_total",
		Help: "api_method_counted_total",
	},
	// 设置标签 系统名称、模块名称、IP、请求类型、URI、返回状态
	[]string{"application", "description", "ip", "method", "uri", "port", "status"},
)

// RequestTime 请求响应时间
var RequestTime = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "api_method_timed_seconds_max",
		Help: "api_method_timed_seconds_max",
	},
	// 设置标签 系统名称、模块名称、IP、请求类型、URI
	[]string{"application", "description", "ip", "method", "port", "uri"},
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type Response struct {
	Code int `json:"code"`
}

package ico

import (
	"bigrule/common/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type IController interface {
	DoHandle(c *gin.Context) *Result
}

func Handler(controller IController) gin.HandlerFunc {
	return func(c *gin.Context) {

		rst := controller.DoHandle(c)

		switch strings.ToLower(rst.Type) {
		case "json":
			c.JSON(http.StatusOK, rst)

		case "string":
			c.String(http.StatusOK, rst.Message)

		case "file":
		}
	}
}

type Result struct {
	Type    string        `json:"-"`
	Status  int           `json:"status"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Logs    []interface{} `json:"-"`
}

func Succ(data interface{}, logs ...interface{}) *Result {
	if data == nil {
		data = map[string]string{}
	}
	return &Result{
		Type:    "json",
		Status:  1,
		Code:    200,
		Message: "ok",
		Data:    data,
		Logs:    logs,
	}
}

func Err(code int, message string, logs ...interface{}) *Result {
	// 1.获取标准errMessage
	errMessage, ok := global.ResponseMessage[code]
	if !ok {
		// 1.1 不存在且无msg
		if message == "" {
			errMessage = "该code不存在"
		} else {
			// 1.2 其他服务信息
			errMessage = message
			message = ""
		}
	}
	// 2.拼接msg
	if message == "" {
		message = errMessage
	} else {
		message = errMessage + "——" + message
	}
	return &Result{
		Type:    "json",
		Status:  0,
		Code:    code,
		Message: message,
		Data:    map[string]string{},
		Logs:    logs,
	}
}

func ErrJwt(code int, message string, logs ...interface{}) *Result {
	return &Result{
		Type:    "json",
		Status:  0,
		Code:    code,
		Message: message,
		Data:    "",
	}
}

func String(message string) *Result {
	return &Result{
		Type:    "string",
		Status:  1,
		Message: message,
	}
}

func File(message string) *Result {
	return &Result{
		Type:    "file",
		Status:  1,
		Message: message,
	}
}

type JFIController interface {
	DoHandle(c *gin.Context) *JFResult
}

func JFHandler(controller JFIController) gin.HandlerFunc {
	return func(c *gin.Context) {
		rst := controller.DoHandle(c)
		c.JSON(http.StatusOK, rst)
	}
}

type JFResult struct {
	Code    int    `json:"resultCode"`
	Message string `json:"resultDesc"`
}

func JFSucc() *JFResult {
	return &JFResult{
		Code:    0,
		Message: "成功",
	}
}

func JFErr(code int, message string) *JFResult {
	return &JFResult{
		Code:    code,
		Message: message,
	}
}

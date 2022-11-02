package jwt

import (
	"bigrule/common/global"
	"bigrule/common/logger"
	config "bigrule/services/authentication/config"
	"bigrule/services/authentication/middleware"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type JwtClaims struct {
	*jwt.StandardClaims
	UserId      int
	UserName    string
	DeptId      int
	RoleIdList  []int
	GroupIdList []int
}

//自定义的token秘钥
var key = []byte("flowpp@flowpp.com")

// Jwt 验证token
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// prometheus 统计请求状态数量 start1
		startTime := time.Now()
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		// end1
		path := strings.Split(c.Request.RequestURI, "?")[0] //过滤是否验证role
		if utils.IsContains(middleware.AJwtNoVerify, path) {
			c.Next()
			getRequest(c, startTime, blw)
			return
		}
		token := c.GetHeader("X-Access-Token") //cookie中拿到token
		if err := AuthToken(token, c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"code":    302,
				"message": err.Error(),
				"data":    "token: " + token,
			})
			c.Abort()
			return
		}
		// wait to release userLock
		c.Next()
		getRequest(c, startTime, blw)
	}
}

func getRequest(c *gin.Context, startTime time.Time, blw *CustomResponseWriter) {
	// prometheus 统计请求状态数量 start2
	if !utils.IsContains(config.PrometheusConfig.UrlList, c.Request.RequestURI) {
		return
	}
	prometheusRes := Response{}
	if err := json.Unmarshal([]byte(blw.body.String()), &prometheusRes); err != nil {
		logger.Error(blw.body.String(), " json解析异常")
		return
	}
	prometheusStatus := "200"
	if prometheusRes.Code != 200 {
		prometheusStatus = "504"
	}
	allLabels := prometheus.Labels{
		"application": "大屏", "description": "用户认证服务", "ip": config.ApplicationConfig.Host, "port": config.ApplicationConfig.Port,
		"method": c.Request.Method, "uri": c.Request.RequestURI, "status": prometheusStatus}
	RequestAllCount.With(allLabels).Inc()
	costTime := time.Since(startTime).Nanoseconds()
	costTimeNano, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(costTime)/1e6), 64)
	labels := prometheus.Labels{
		"application": "大屏", "description": "用户认证服务", "ip": config.ApplicationConfig.Host, "port": config.ApplicationConfig.Port,
		"method": c.Request.Method, "uri": c.Request.RequestURI}
	RequestTime.With(labels).Set(costTimeNano)
	// end2
}

// ManageAuthToken 其他服务验证token
func ManageAuthToken(path string, c *gin.Context) (code int, err error) {
	if utils.IsContains(middleware.JwtNoVerify, path) {
		return 200, nil
	}
	token := c.GetHeader("X-Access-Token") //cookie中拿到token
	if err = AuthToken(token, c); err != nil {
		return 302, err
	}
	return 200, nil
}

// GenToken 生成token
func GenToken(up map[string]interface{}, dt int64) (string, error) {
	claims := JwtClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: dt,
		},
		UserId:   up["user_id"].(int),
		UserName: up["user_name"].(string),
		DeptId:   up["dept_id"].(int),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// ParseJwt 解析token
func ParseJwt(token string) (*JwtClaims, error) {
	var jClaim = &JwtClaims{}
	getToken, err := jwt.ParseWithClaims(token, jClaim, func(*jwt.Token) (interface{}, error) { return key, nil })
	if err != nil {
		fmt.Println("无法处理这个token", err)
		return jClaim, errors.New("token 错误")
	}
	if !getToken.Valid { //服务端验证token是否有效
		if ve, ok := err.(*jwt.ValidationError); ok { //官方写法招抄就行
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("错误的token")
				return jClaim, errors.New("cookie异常：token 无效")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("token过期或未启用")
				return jClaim, errors.New("cookie异常：token 失效")
			} else {
				fmt.Println("无法处理这个token", err)
				return jClaim, errors.New("cookie异常：token 错误")
			}
		}
	}
	return jClaim, nil
}

// 检验mysql token
func mysqlToken(userId int, token string) (err error) {
	tokenT := users.TokenTable{UserId: userId}
	tokenInfo, err := tokenT.Query(global.DBMysql)
	if err == nil && tokenInfo.Token == token && tokenInfo.DueDate >= time.Now().Unix() {
		return nil
	}
	return errors.New("token 失效")
}

func AuthToken(token string, c *gin.Context) (err error) {
	// 1.是否存在
	if token == "" || token == "undefined" {
		return errors.New("cookie异常：token 不存在")
	}
	token = strings.Split(token, ";")[0]
	// 2.是否正常使用
	reqMsg, err := ParseJwt(token)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return err
	}
	// 3.单点登录
	if err = mysqlToken(reqMsg.UserId, token); err != nil {
		logger.Error(err)
		c.Abort()
		return err
	}
	c.Set("user_id", reqMsg.UserId)
	c.Set("user_name", reqMsg.UserName)
	c.Set("dept_id", reqMsg.DeptId)
	c.Set("token", token)
	return
}

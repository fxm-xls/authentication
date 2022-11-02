package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/config"
	"bigrule/services/authentication/middleware/jwt"
	"bigrule/services/authentication/model/users"
	"bigrule/services/authentication/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type TokenLogin struct {
	Token string `json:"token" binding:"required"`
	Keys  string `json:"keys" binding:"required"`
}

type TokenLoginRes struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (This TokenLogin) DoHandle(c *gin.Context) *ico.Result {
	logger.Info("CheckUrl:", config.JFConfig.CheckUrl)
	logger.Info("config.JFConfig.UserUrl:", config.JFConfig.UserUrl)
	logger.Info("config.JFConfig.CsrUrl:", config.JFConfig.CsrUrl)
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "")
	}
	logger.Infof("经分用户登陆:")
	logger.Infof("This.Keys:", This.Keys)
	logger.Infof("This.Token:", This.Token)
	// 1.验证token
	success, err := CheckToken(This.Token)
	if err != nil && success {
		logger.Error(success, err)
		return ico.Err(2001, err.Error())
	}
	// 2.解析keys
	operId, idType, err := This.SplitKeys()
	if err != nil {
		logger.Error(err)
		return ico.Err(2001, err.Error())
	}
	// 3.获取用户信息
	logger.Info("3.获取用户信息")
	userInfo, err := GetUser(operId, idType)
	if err != nil {
		logger.Error(err)
		return ico.Err(2001, err.Error())
	}
	// 4.Csr登录
	logger.Info("4.Csr登录")
	res, err := This.GetCsrToken(userInfo.Body.OperId, userInfo.Body.OperName)
	if err != nil {
		logger.Error(err)
		return ico.Err(2001, err.Error())
	}
	return ico.Succ(res)
}

func (This TokenLogin) SplitKeys() (operId, idType string, err error) {
	keysBytes, err := base64.StdEncoding.DecodeString(This.Keys)
	if err != nil {
		logger.Error(This.Keys, err)
		return operId, idType, errors.New("解码keys失败")
	}
	logger.Info("string(keysBytes):", string(keysBytes))
	keysList := strings.Split(string(keysBytes), "&")
	for _, idStr := range keysList {
		logger.Info("idStr:", idStr)
		idList := strings.Split(idStr, "=")
		if len(idList) < 1 {
			logger.Info(idList, "keys切割异常")
		}
		operId = idList[1]
		if idList[0] == "operId" {
			idType = "operId"
		} else if idList[0] == "msisdn" {
			idType = "msisdn"
		}
	}
	return
}

type TokenResponse struct {
	Code    string `json:"resultCode"`
	Message string `json:"resultDesc"`
}

func CheckToken(token string) (success bool, err error) {
	success = true
	var tokenRes TokenResponse
	//urlToken := fmt.Sprintf("http://10.32.40.78:30012/sso_server/api/checkToken.action")
	//CsrUrl := "10.32.233.56:8081"
	urlToken := config.JFConfig.CheckUrl
	tokenPars := map[string]string{"targetIP": config.JFConfig.CsrUrl, "token": token, "sysId": "H00000000091"}
	headers := map[string]string{}
	respToken, err := utils.PostUrl(tokenPars, urlToken, headers)
	if err != nil {
		err = errors.New("token检测失败")
		logger.Error(err)
		return
	}
	err = json.Unmarshal(respToken, &tokenRes)
	if err != nil {
		err = errors.New("token检测失败")
		logger.Error(err)
		return
	}
	if tokenRes.Code != "0" {
		return false, errors.New(tokenRes.Message)
	}
	return
}

type UserResponse struct {
	Code    string `json:"resultCode"`
	Message string `json:"resultDesc"`
	Body    Body   `json:"body"`
}

type Body struct {
	OperId     string `json:"operId"`
	OperName   string `json:"opername"`
	Msisdn     string `json:"msisdn"`
	Email      string `json:"email"`
	Sex        string `json:"sex"`
	RegionId   string `json:"regionId"`
	RegionName string `json:"regionName"`
	OrgId      string `json:"orgId"`
	OrgName    string `json:"orgName"`
}

func GetUser(operId, idType string) (user UserResponse, err error) {
	//urlUser := fmt.Sprintf("http://10.32.40.78:30011/singlePoint/getUserInfo.action")
	urlUser := config.JFConfig.UserUrl
	tokenPars := map[string]string{"sysId": "H00000000091", "operId": operId, "msisdn": operId}
	if idType == "operId" {
		tokenPars["type"] = "1"
	} else if idType == "msisdn" {
		tokenPars["type"] = "2"
	}
	logger.Info("tokenPars:", tokenPars)
	headers := map[string]string{}
	respUser, err := utils.PostUrl(tokenPars, urlUser, headers)
	if err != nil {
		err = errors.New("经分用户查询失败")
		logger.Error(err)
		return
	}
	logger.Info("respUser:", string(respUser))
	if err = json.Unmarshal(respUser, &user); err != nil {
		err = errors.New("经分用户查询失败")
		logger.Error(err)
		return
	}
	if user.Code != "0" {
		logger.Error(user)
		return user, errors.New(user.Message)
	}
	return
}

func (This TokenLogin) GetCsrToken(account, userName string) (Res TokenLoginRes, err error) {
	// 1.1验证用户账户
	logger.Info("1.1验证用户账户")
	userT := &users.UserTable{}
	user, err := userT.QueryByFilter(global.DBMysql, map[string]interface{}{"account": account})
	if err != nil {
		err = errors.New("operId错误")
		logger.Error(err)
		return
	}
	// 1.2.验证用户名
	logger.Info("1.2.验证用户名")
	if userName != user.UserName {
		err = errors.New("用户名错误")
		logger.Error(err)
		return
	}
	// 2.设置token
	logger.Info("2.设置token")
	userMsg := make(map[string]interface{})
	userMsg["user_id"] = user.UserId
	userMsg["user_name"] = user.UserName
	userMsg["dept_id"] = 0
	// 2.1.token有效时间
	logger.Info("2.1.token有效时间")
	expireTimes := config.CookieConfig.MysqlTime * time.Second
	expires := time.Now().Add(expireTimes).Unix()
	token, err := jwt.GenToken(userMsg, expires)
	if err != nil {
		err = errors.New("登陆失败")
		logger.Error(err)
		return
	}
	tx := global.DBMysql
	// 2.2.mysql 存入token
	logger.Info("2.2.mysql 存入token")
	if err = users.SetToken(tx, user.UserId, token, expires); err != nil {
		err = errors.New("登陆失败")
		logger.Error(err)
		return
	}
	Res.UserId = user.UserId
	Res.Token = token
	return
}

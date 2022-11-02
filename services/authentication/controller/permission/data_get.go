package permission

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/users"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DataGet struct {
	ServiceName string `json:"service_name" binding:"required"`
}

type DataGetRes struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	Operation string `json:"operation"`
}

func (This DataGet) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "")
	}
	userId := c.GetInt("user_id")
	logger.Infof("数据查询 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取单个服务
	res, code, err := This.GetData(userId)
	if err != nil {
		return ico.Err(code, err.Error())
	}
	return ico.Succ(res)
}

func (This DataGet) JudgeInfo(userId int) (code int, err error) {
	// 1.检测用户是否存在
	userT := users.UserTable{}
	if _, err = userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": userId}); err != nil {
		logger.Error("用户不存在", userId)
		return 2144, errors.New(fmt.Sprint("userId：", userId))
	}
	return
}

func (This DataGet) GetData(userId int) (res []DataGetRes, code int, err error) {
	// 1.获取数据
	dataV := entity.DataView{}
	whereMap := map[string]interface{}{"user_id": userId, "service_name": This.ServiceName}
	dataList, err := dataV.QueryListByFilter(global.DBMysql, whereMap)
	if err != nil {
		logger.Info(whereMap, "数据查询失败")
		return res, 2162, errors.New(fmt.Sprintf("user_id:%d,service_name:%s", userId, This.ServiceName))
	}
	for _, data := range dataList {
		res = append(res, DataGetRes{DataId: data.DataId, DataType: data.DataType, Operation: data.Operation})
	}
	return
}

package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserDataQuery struct {
	ServiceName string `json:"service_name" binding:"required"`
	UserId      int    `json:"user_id"      binding:"required"`
}

type UserDataQueryRes struct {
	DataId     int    `json:"data_id"`
	DataName   string `json:"data_name"`
	DataType   string `json:"data_type"`
	DataDesc   string `json:"data_desc"`
	Operation  string `json:"operation"`
	Status     int    `json:"status"`
	BindStatus int    `json:"bind_status"`
}

func (This UserDataQuery) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("用户查询_数据 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.获取用户
	res, code, err := This.GetUsers()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.空数组
	if len(res) == 0 {
		res = []UserDataQueryRes{}
	}
	return ico.Succ(res)
}

func (This UserDataQuery) JudgeInfo(userId int) (code int, err error) {
	// 1.检测账户是否存在
	userT := users.UserTable{}
	if _, err = userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.UserId}); err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	// 2.检测是否是服务管理员
	//if code, err = public.JudgeManager(userId, This.ServiceName); err != nil {
	//	logger.Error(err.Error())
	//	return code, err
	//}
	return
}

func (This UserDataQuery) GetUsers() (res []UserDataQueryRes, code int, err error) {
	// 0.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("该服务不存在")
		return res, 2104, errors.New("")
	}
	// 1.获取所有数据
	entityT := entity.EntityTable{}
	entityList, err := entityT.QueryListByFilter(global.DBMysql, map[string]interface{}{"service_id": serviceId})
	if err != nil {
		logger.Error("角色查询失败")
		return res, 2112, errors.New("")
	}
	// 2.获取该用户数据
	userEntityT := users.UserEntityTable{}
	userEntityList, err := userEntityT.QueryListByFilter(global.DBMysql, map[string]interface{}{"user_id": This.UserId})
	if err != nil {
		logger.Error("用户与数据联系查询失败")
		return res, 2143, errors.New("")
	}
	// 3.获取所有用户数据
	userEntityAll, err := userEntityT.QueryListByFilter(global.DBMysql, map[string]interface{}{})
	if err != nil {
		logger.Error("用户与数据联系查询失败")
		return res, 2143, errors.New("")
	}
	// 4.配置
	for _, entityInfo := range entityList {
		dataInfo := UserDataQueryRes{
			DataId: entityInfo.DataId, DataName: entityInfo.DataName, DataDesc: entityInfo.DataDesc,
			DataType: entityInfo.DataType, Operation: entityInfo.Operation,
		}
		for _, userEntity := range userEntityList {
			if userEntity.EntityId == entityInfo.EntityId {
				dataInfo.Status = 1
				break
			}
		}
		for _, userEntity := range userEntityAll {
			if userEntity.EntityId == entityInfo.EntityId {
				dataInfo.BindStatus = 1
				break
			}
		}
		res = append(res, dataInfo)
	}
	return
}

package users

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDataUpdate struct {
	ServiceName string               `json:"service_name" binding:"required"`
	UserId      int                  `json:"user_id"      binding:"required"`
	DataList    []UserDataUpdateInfo `json:"data_list"    binding:"required"`
}

type UserDataUpdateInfo struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	Operation string `json:"operation"`
}

func (This UserDataUpdate) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("修改用户数据 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.校验
	if code, err := This.JudgeInfo(userId); err != nil {
		return ico.Err(code, err.Error())
	}
	// 2.修改用户数据联系
	tx := global.DBMysql.Begin()
	if code, err := This.UpdateUserData(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("用户数据修改成功")
}

func (This UserDataUpdate) JudgeInfo(userId int) (code int, err error) {
	// 1.检测账户是否存在
	userT := users.UserTable{}
	if _, err = userT.QueryByFilter(global.DBMysql, map[string]interface{}{"id": This.UserId}); err != nil {
		logger.Error("用户不存在")
		return 2144, errors.New("")
	}
	// 2.检测是否是服务管理员
	if err = model.JudgeManager(userId, global.CsrServiceName); err != nil {
		logger.Error(err.Error())
		return code, err
	}
	return
}

func (This UserDataUpdate) UpdateUserData(tx *gorm.DB) (code int, err error) {
	// 1.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("该服务不存在")
		return 2104, errors.New("")
	}
	// 2.获取服务下所有数据
	entityT := entity.EntityTable{}
	entityList, err := entityT.QueryListByFilter(global.DBMysql, map[string]interface{}{"service_id": serviceId})
	if err != nil {
		logger.Error("数据查询失败")
		return 2162, errors.New("")
	}
	// 3.获取entityId
	var userEntityList []users.UserEntityTable
	for _, data := range This.DataList {
		for _, entityInfo := range entityList {
			if data.Operation == "1" {
				continue
			}
			if data.DataId == entityInfo.DataId && data.DataType == entityInfo.DataType && data.Operation == entityInfo.Operation {
				userEntityList = append(userEntityList, users.UserEntityTable{UserId: This.UserId, EntityId: entityInfo.EntityId})
				break
			}
		}
	}
	// 4.删除用户数据联系
	userEntityT := users.UserEntityTable{}
	if err = userEntityT.DeleteByUserIds(tx, []int{This.UserId}); err != nil {
		logger.Error("用户与数据联系删除失败")
		return 2137, errors.New("")
	}
	// 5.新增用户数据联系
	if len(userEntityList) > 0 {
		if err = userEntityT.InsertMany(tx, userEntityList); err != nil {
			logger.Error("用户与数据联系修改失败")
			return 2147, errors.New("")
		}
	}
	return
}

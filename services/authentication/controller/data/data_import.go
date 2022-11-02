package data

import (
	"bigrule/common/global"
	"bigrule/common/ico"
	"bigrule/common/logger"
	"bigrule/services/authentication/model/entity"
	"bigrule/services/authentication/model/services"
	"bigrule/services/authentication/model/users"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DataImport struct {
	ServiceName    string     `json:"service_name"   binding:"required"`
	AddDataList    []DataInfo `json:"add_list"       `
	UpdateDataList []DataInfo `json:"update_list"    `
	DelDataList    []DataInfo `json:"delete_list"    `
	ServiceId      int
	EntityT        entity.EntityTable
}

type DataInfo struct {
	DataId    int    `json:"data_id"`
	DataType  string `json:"data_type"`
	DataName  string `json:"data_name"`
	DataDesc  string `json:"data_desc"`
	Operation string `json:"operation"`
}

func (This DataImport) DoHandle(c *gin.Context) *ico.Result {
	if err := c.ShouldBindJSON(&This); err != nil {
		return ico.Err(2099, "", err.Error())
	}
	userId := c.GetInt("user_id")
	logger.Infof("数据导入 user: id %d, name %s", userId, c.GetString("user_name"))
	// 1.获取服务id
	serviceId, err := services.GetServiceIdByName(This.ServiceName)
	if err != nil {
		logger.Error("服务不存在 ", This.ServiceName)
		return ico.Err(2104, This.ServiceName)
	}
	This.ServiceId = serviceId
	// 2.校验
	This.EntityT = entity.EntityTable{}
	updateIdList, dalIdList, code, err := This.JudgeInfo()
	if err != nil {
		return ico.Err(code, err.Error())
	}
	// 3.删除
	tx := global.DBMysql.Begin()
	if code, err = This.DelData(tx, dalIdList); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	// 4.增加
	if code, err = This.AddData(tx); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	// 5.修改
	if code, err = This.UpdateData(tx, updateIdList); err != nil {
		tx.Rollback()
		return ico.Err(code, err.Error())
	}
	tx.Commit()
	return ico.Succ("导入成功")
}

func (This DataImport) JudgeInfo() (updateIdList, dalIdList []int, code int, err error) {
	// 1.增加数据是否存在
	for _, data := range This.AddDataList {
		whereMap := map[string]interface{}{
			"data_id": data.DataId, "data_type": data.DataType, "operation": data.Operation, "service_id": This.ServiceId,
		}
		if _, err = This.EntityT.QueryByFilter(global.DBMysql, whereMap); err == nil {
			message := fmt.Sprintf("data_id:%d,data_type:%s,operation:%s", data.DataId, data.DataType, data.Operation)
			logger.Error("数据已存在 ", message)
			return updateIdList, dalIdList, 2165, errors.New(message)
		}
	}
	err = nil
	// 2.修改、删除数据是否存在
	for _, data := range This.UpdateDataList {
		whereMap := map[string]interface{}{
			"data_id": data.DataId, "data_type": data.DataType, "operation": data.Operation, "service_id": This.ServiceId,
		}
		entityInfo, err := This.EntityT.QueryByFilter(global.DBMysql, whereMap)
		if err != nil {
			message := fmt.Sprintf("data_id:%d,data_type:%s,operation:%s", data.DataId, data.DataType, data.Operation)
			logger.Error("修改数据不存在 ", message)
			return updateIdList, dalIdList, 2166, errors.New(message)
		}
		updateIdList = append(updateIdList, entityInfo.EntityId)
	}
	for _, data := range This.DelDataList {
		whereMap := map[string]interface{}{
			"data_id": data.DataId, "data_type": data.DataType, "operation": data.Operation, "service_id": This.ServiceId,
		}
		entityInfo, err := This.EntityT.QueryByFilter(global.DBMysql, whereMap)
		if err != nil {
			message := fmt.Sprintf("data_id:%d,data_type:%s,operation:%s", data.DataId, data.DataType, data.Operation)
			logger.Error("删除数据不存在 ", message)
			return updateIdList, dalIdList, 2166, errors.New(message)
		}
		dalIdList = append(dalIdList, entityInfo.EntityId)
	}
	return
}

func (This DataImport) DelData(tx *gorm.DB, dalIdList []int) (code int, err error) {
	if len(This.DelDataList) == 0 {
		return
	}
	// 1.删除数据
	if err = This.EntityT.DeleteByIds(tx, dalIdList); err != nil {
		logger.Error("数据删除失败")
		return 2163, errors.New("")
	}
	// 2.删除数据与用户联系
	userEntityT := users.UserEntityTable{}
	if err = userEntityT.DeleteByEntityIds(tx, dalIdList); err != nil {
		logger.Error("用户与数据联系删除失败")
		return 2137, errors.New("")
	}
	return
}

func (This DataImport) AddData(tx *gorm.DB) (code int, err error) {
	if len(This.AddDataList) == 0 {
		return
	}
	var entityList []entity.EntityTable
	for _, data := range This.AddDataList {
		entityList = append(entityList, entity.EntityTable{
			DataId: data.DataId, DataType: data.DataType, DataName: data.DataName,
			DataDesc: data.DataDesc, Operation: data.Operation, ServiceId: This.ServiceId,
		})
	}
	if len(entityList) > 0 {
		if err = This.EntityT.InsertMany(tx, entityList); err != nil {
			logger.Error("数据导入失败")
			return 2161, errors.New("")
		}
	}
	return
}

func (This DataImport) UpdateData(tx *gorm.DB, updateIdList []int) (code int, err error) {
	if len(This.UpdateDataList) == 0 {
		return
	}
	for i, data := range This.UpdateDataList {
		entityT := entity.EntityTable{
			EntityId: updateIdList[i], DataName: data.DataName, DataDesc: data.DataDesc,
		}
		if err = entityT.UpdateByStruct(tx); err != nil {
			message := fmt.Sprintf("data_id:%d,data_type:%s,operation:%s", data.DataId, data.DataType, data.Operation)
			logger.Error("数据修改失败 ", message)
			return 2164, errors.New(message)
		}
	}
	return
}

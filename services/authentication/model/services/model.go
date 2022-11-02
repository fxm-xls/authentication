package services

import (
	"bigrule/common/global"
)

// GetServiceIdByName 根据服务名称获取服务id
func GetServiceIdByName(serviceName string) (serviceId int, err error) {
	serviceT := ServiceTable{}
	service, err := serviceT.QueryByFilter(global.DBMysql, map[string]interface{}{"name": serviceName})
	if err != nil {
		return
	}
	serviceId = service.ServiceId
	return
}

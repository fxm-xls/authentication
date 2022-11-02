package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"gorm.io/gorm"
)

const (
	ProjectName = "BigRule"
	Version     = "v3"
	// DefaultUserInt 内置用户 default字段默认值
	DefaultUserInt = 1
	// DefaultJFUserInt 经分内置用户 default字段默认值
	DefaultJFUserInt = 2
	// DefaultManagerInt 内置角色 manager字段默认值
	DefaultManagerInt = 1
	// DefaultServiceId 内置服务id
	DefaultServiceId = 1
	// DefaultServiceName 内置服务名称
	DefaultServiceName = "authentication-service"
	CsrServiceName     = "flowcsr-service"
	RepoServiceName    = "repo-service"
	// AuthManagerId 权限服务管理员id
	AuthManagerId = 1
	// DefaultSuffix 管理员默认后缀
	DefaultSuffix = "manager"
	// DefaultLimitNum 角色默认限制数量
	DefaultLimitNum = -1
)

var (
	GinEngine      *gin.Engine
	CasbinEnforcer *casbin.SyncedEnforcer
	EtcdReg        registry.Registry
	DBMysql        *gorm.DB
	// CsrServiceIds csr服务默认id列表
	CsrServiceIds = []int{6, 7, 8}
	// ResponseMessage 接口返回信息
	ResponseMessage map[int]string
)

var Logo = []byte{
	10, 32, 32, 32, 32, 47, 47, 32, 32, 32, 41, 32, 41, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 47, 47, 32, 32, 32, 41, 32, 41, 10, 32, 32, 32, 47, 47,
	95, 95, 95, 47, 32, 47, 32, 32, 32, 32, 32, 40, 32, 41, 32, 32, 32, 32, 32, 95, 95, 95, 32, 32, 32,
	32, 32, 32, 32, 47, 47, 95, 95, 95, 47, 32, 47, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 47, 47, 32, 32, 32, 32, 32, 32, 95, 95, 95, 10, 32, 32, 47, 32, 95, 95, 32, 32, 40, 32,
	32, 32, 32, 32, 32, 47, 32, 47, 32, 32, 32, 32, 47, 47, 32, 32, 32, 41, 32, 41, 32, 32, 32, 47, 32,
	95, 95, 95, 32, 40, 32, 32, 32, 32, 32, 32, 47, 47, 32, 32, 32, 47, 32, 47, 32, 32, 32, 47, 47, 32,
	32, 32, 32, 32, 47, 47, 95, 95, 95, 41, 32, 41, 10, 32, 47, 47, 32, 32, 32, 32, 41, 32, 41, 32, 32,
	32, 32, 47, 32, 47, 32, 32, 32, 32, 40, 40, 95, 95, 95, 47, 32, 47, 32, 32, 32, 47, 47, 32, 32, 32,
	124, 32, 124, 32, 32, 32, 32, 32, 47, 47, 32, 32, 32, 47, 32, 47, 32, 32, 32, 47, 47, 32, 32, 32,
	32, 32, 47, 47, 10, 47, 47, 95, 95, 95, 95, 47, 32, 47, 32, 32, 32, 32, 47, 32, 47, 32, 32, 32, 32,
	32, 32, 47, 47, 95, 95, 32, 32, 32, 32, 32, 47, 47, 32, 32, 32, 32, 124, 32, 124, 32, 32, 32, 32,
	40, 40, 95, 95, 95, 40, 32, 40, 32, 32, 32, 47, 47, 32, 32, 32, 32, 32, 40, 40, 95, 95, 95, 95, 10,
}

package mycasbin

import (
	"bigrule/common/global"
	"bigrule/common/logger"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch4(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*") || r.sub == "1"
`

func Setup() {
	Apter, err := gormAdapter.NewAdapterByDB(global.DBMysql)
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewSyncedEnforcer(m, Apter)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
	global.CasbinEnforcer = e
}

func Casbin() *casbin.SyncedEnforcer {
	return global.CasbinEnforcer
}

func LoadPolicy() (*casbin.SyncedEnforcer, error) {
	if err := global.CasbinEnforcer.LoadPolicy(); err == nil {
		return global.CasbinEnforcer, err
	} else {
		logger.Infof("casbin rbac_model or policy init error, message: %v \r\n", err.Error())
		return nil, err
	}
}

// AuthCheckRole 检查权限
func AuthCheckRole(roleIds []int, urlPath, method string) (err error) {
	e := Casbin()
	for _, roleId := range roleIds {
		res, err := e.Enforce(fmt.Sprint(roleId), urlPath, method)
		if err == nil && res {
			return nil
		}
	}
	logger.Info("校验:", urlPath)
	return errors.New("权限不足")
}

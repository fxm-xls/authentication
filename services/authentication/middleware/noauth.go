package middleware

import "bigrule/common/global"

// JwtNoVerify 该路由下免登录不校验
var JwtNoVerify = []string{
	"/authentication-service/v2/users/login", "/authentication-service/" + global.Version + "/users/login",
	"/authentication-service/v2/users/jf-service/login",
	"/authentication-service/v2/users/jf-service/add",
}
var AJwtNoVerify = []string{
	"/v2/users/login", "/" + global.Version + "/users/login",
	"/v2/users/jf-service/login", "/v2/services/register", "/v2/users/jf-service/add",
}

// CasbinNoVerify 该路由下不校验
var CasbinNoVerify = []string{
	"/authentication-service/v3/users/token/query", "/authentication-service/v3/users/logout",
	"/authentication-service/v3/users/query", "/authentication-service/v3/users/update",
	"/authentication-service/v3/users/data/query",
	"/authentication-service/v3/data/user/query",
	"/authentication-service/v3/permission/data/get",
	"/authentication-service/v3/user/manager/query",

	"/authentication-service/v2/users/token/query", "/authentication-service/v2/users/logout",
	"/authentication-service/v2/users/query", "/authentication-service/v2/users/update",
	"/authentication-service/v2/users/data/query",
	"/authentication-service/v2/data/user/query",
	"/authentication-service/v2/permission/data/get",
	"/authentication-service/v2/user/manager/query",
}

// ManagerCasbin 该路由下管理员不校验
var ManagerCasbin = []string{
	////服务模块
	//"/authentication-service/v2/services/delete", "/authentication-service/v2/services/update",
	//// 用户模块
	//"/authentication-service/v2/users/data/update", "/authentication-service/v2/users/roles/query",
	//// 角色模块
	//"/authentication-service/v2/roles/add", "/authentication-service/v2/roles/delete",
	//"/authentication-service/v2/roles/query", "/authentication-service/v2/roles/update",
	//"/authentication-service/v2/roles/interface/query", "/authentication-service/v2/roles/interface/update",
	//// 接口模块
	//"/authentication-service/v2/interfaces/query", "/authentication-service/v2/interfaces/import",
	//// 数据模块
	//"/authentication-service/v2/data/query", "/authentication-service/v2/data/import",
	//// v2模块
	//"/authentication-service/v2/user/permission/update", "/authentication-service/v2/user/permission/query",
}

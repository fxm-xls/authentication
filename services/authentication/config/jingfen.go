package config

import "github.com/spf13/viper"

type JF struct {
	CheckUrl string
	UserUrl  string
	CsrUrl   string
	RoleId   int
	DeptId   int
	IPS      []string
}

func InitJF(cfg *viper.Viper) *JF {
	return &JF{
		CheckUrl: cfg.GetString("checkUrl"),
		UserUrl:  cfg.GetString("userUrl"),
		CsrUrl:   cfg.GetString("csrUrl"),
		RoleId:   cfg.GetInt("roleId"),
		DeptId:   cfg.GetInt("deptId"),
		IPS:      cfg.GetStringSlice("ips"),
	}
}

var JFConfig = new(JF)

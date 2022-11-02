package config

import (
	"bigrule/common/etcd"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// database config
var cfgDatabase *viper.Viper

// application config
var cfgApplication *viper.Viper

// log config
var cfgLogger *viper.Viper

// etcd config
var cfgEtcd *viper.Viper

// cookie config
var cfgCookie *viper.Viper

// 经分 config
var cfgJF *viper.Viper

// web static config
var cfgWebStatic *viper.Viper

// prometheus config
var cfgPrometheus *viper.Viper

// cfgResMsg config
var cfgResMsg *viper.Viper

// Setup config
func Setup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
	//DbMysql
	cfgDatabase = viper.Sub("bigrule.authentication.dbmysql")
	if cfgDatabase == nil {
		panic("No found bigrule.authentication.dbmysql in the configuration")
	}
	DbMysqlConfig = InitDbMysql(cfgDatabase)
	//application
	cfgApplication = viper.Sub("bigrule.authentication.application")
	if cfgApplication == nil {
		panic("No found bigrule.authentication.application in the configuration")
	}
	ApplicationConfig = InitApplication(cfgApplication)
	//logger
	cfgLogger = viper.Sub("bigrule.authentication.logger")
	if cfgLogger == nil {
		panic("No found bigrule.authentication.logger in the configuration")
	}
	LoggerConfig = InitLogger(cfgLogger)
	//cookie
	cfgCookie = viper.Sub("bigrule.authentication.cookie")
	if cfgCookie == nil {
		panic("No found bigrule.authentication.cookie in the configuration")
	}
	CookieConfig = InitCookie(cfgCookie)
	//etcd
	cfgEtcd = viper.Sub("bigrule.etcd-service")
	if cfgEtcd == nil {
		panic("No found bigrule.etcd-service in the configuration")
	}
	etcd.EtcdConfig = etcd.InitEtcd(cfgEtcd)
	//经分
	cfgJF = viper.Sub("bigrule.authentication.jf")
	if cfgJF == nil {
		panic("No found bigrule.authentication.jf in the configuration")
	}
	JFConfig = InitJF(cfgJF)
	//web static
	cfgWebStatic = viper.Sub("bigrule.authentication.web-static")
	if cfgWebStatic == nil {
		panic("No found bigrule.authentication.web-static in the configuration")
	}
	WebStaticConfig = InitWebStatic(cfgWebStatic)
	// Prometheus
	cfgPrometheus = viper.Sub("bigrule.authentication.prometheus")
	if cfgPrometheus == nil {
		panic("Not found bigrule.authentication.prometheus in the configuration")
	}
	PrometheusConfig = InitPrometheus(cfgPrometheus)
	// cfgResMsg
	cfgResMsg = viper.Sub("bigrule.authentication.response-message")
	if cfgResMsg == nil {
		panic("Not found bigrule.authentication.response-message in the configuration")
	}
	ResMsgConfig = InitResMsg(cfgResMsg)
	//......
}

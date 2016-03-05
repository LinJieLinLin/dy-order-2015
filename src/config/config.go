package config

import (
	"github.com/Centny/gwf/util"
	"github.com/Centny/gwf/log"
)

//Orders configure.
var Cfg *util.Fcfg = util.NewFcfg3()

func DbConfig() string {
	return Cfg.Val("DB_CONFIG")
}
func TestDbConfig() string {
	return Cfg.Val("TEST_DB_CONFIG")
}
func ServerPort() string {
	return Cfg.Val("SERVER_PORT")
}
func LogPath() string {
	return Cfg.Val("LOG_PATH")
}
func LogFileName() string {
	return Cfg.Val("LOG_FILE_NAME")
}
func AlipayNotifyHost() string {
	return Cfg.Val("ALIPAY_NOTIFY_HOST")
}
//get the SSO logout URL.
func SsoLogoutUrl() string{
	return Cfg.Val("SSO_LOGOUT_URL")
}

//show the configure.
func ShowConf() {
	log.I(`Server conf(
		DB_CONFIG:%v,
		TEST_DB_CONFIG:%v,
		SERVER_PORT:%v,
		LOG_PATH:%v,
		LOG_FILE_NAME:%v,
		ALIPAY_NOTIFY_HOST:%v,
		SSO_LOGOUT_URL:%v,
		END--------------
		)`,
		DbConfig(), TestDbConfig(), ServerPort(), LogPath(), LogFileName(), AlipayNotifyHost(), SsoLogoutUrl())
}
func init() {
	var cfg string = ""
	cfg = "conf/orders.properties"
	err := Cfg.InitWithFilePath(cfg)
	if err != nil {
		log.E("read config failed --:%v", err)
		return
	}
}
//get config data for test
func TestInit(arg_cfg string) {
	log.D("set db config for test")
	cfg:= "../../../conf/orders.properties"
	if ""!=arg_cfg{
		cfg = arg_cfg
	}
	err := Cfg.InitWithFilePath(cfg)
	if err != nil {
		log.E("read config failed --:%v", err)
		return
	}
}

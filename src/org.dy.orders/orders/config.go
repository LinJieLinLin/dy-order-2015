package orders

import (
	"github.com/Centny/gwf/util"
	"github.com/Centny/gwf/log"
)

//Order configure.
var Cfg *util.Fcfg = util.NewFcfg3()

func AddToOrdersUrl() string {
	return Cfg.Val("ORD_ADD_CART")
}

func init() {
	var cfg string = ""
	cfg = "conf/org.dy.orders.properties"
	err := Cfg.InitWithFilePath(cfg)
	if err != nil {
		log.E("read config failed --:%v", err)
		return
	}
}

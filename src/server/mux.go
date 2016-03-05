package server
import (
	"org.cny.uas/sso"
	"github.com/Centny/gwf/routing"
	"org.cny.uap/conf"
	"org.cny.uap/uap"
	"org.cny.uas/usr"
	"org.cny.uap/sync"
	"github.com/Centny/gwf/log"
	ucs "org.cny.uas/common"
	"net/http"
	"databaseConn"
	C "common"
	"orders/cart"
	"orders/orders"
	O "org.dy.orders/orders"
	"fmt"
	"alipay"
	"github.com/Centny/gwf/routing/filter"
	"config"
)


func RouteConfig() {
	sso.ShowLog=true
	sb := routing.NewSrvSessionBuilder("", "/", "orders", 6000000, 1000)
	mux := routing.NewSessionMux("", sb)
	cors := filter.NewCORS()
	cors.AddSite("*")
	mux.HFilter("^/.*$", cors)
	initUap()
	INT, _ := routing.NewJsonINT("conf/i18n")
	INT.Default = "en"
	mux.INT = INT
//	mux.ShowLog = true
	afc := sso.NewAuthFilter(
		conf.SsoLoginUrl(),
		conf.SsoAuthUrl(),
		"",
	)
	af := sso.NewAuthFilter(
		conf.SsoLoginUrl(),
		conf.SsoAuthUrl(),
		"",
	)
	//afc return 301
	afc.M = "C"
	//logout
	mux.HFilterFunc("^/quitSystem(\\?.*)?$", uap.ClsUsrStmtFilter)
	mux.HFilterFunc("^/quitSystem(\\?.*)?$", sso.ClsAuthFilter)
	mux.HFilterFunc("^/quitSystem(\\?.*)?$", C.Quit)
	mux.HFilter("^/quitSystem(\\?.*)?$", sso.NewRedirect3(config.SsoLogoutUrl(), "http://%s/"))

	//html Filter
	mux.HFilter("/cart.html", af)
	mux.HFilter("/orders.html", af)
	mux.HFilter("/pay.html", af)
	mux.HFilter("/viewOrders.html", af)
	mux.HFilter("^/l/alipay(\\?.*)?$", af)

	//test
	mux.HFunc("/l/test", test)
	//"org.dy.orders/orders"
	mux.HFilterFunc("/aip/test/addToCart", O.AddToCartTest)
	mux.HFilterFunc("/aip/addToCart", O.AddToCart)
	mux.HFilterFunc("/api/EditO_USR", O.EditO_USR)
	mux.HFilterFunc("/api/getUrl", O.GetUrl)
	//"org.dy.orders/orders"
	//testEnd

	mux.HFilter("^/l/.*$", afc)
	//sync user
	mux.HFilter("^/l/.*$", sync.NewSyncUsrStmtFilter())
	//set session
	mux.HFilterFunc("^/l/.*$", C.SessionFilter)

	muxOrder(mux)

	mux.Handler("/", http.FileServer(http.Dir("www/")))
	http.Handle("/", mux)


}
func muxOrder(mux *routing.SessionMux) {
	mux.HFunc("^/l/addToCart(\\?.*)?$", cart.AddToCart)
	mux.HFunc("^/l/getCart(\\?.*)?$", cart.GetCart)
	mux.HFunc("^/l/editCart(\\?.*)?$", cart.EditCart)

	mux.HFunc("^/l/submitOrder(\\?.*)?$", orders.SubmitOrders)
	mux.HFunc("^/l/getOrdersData(\\?.*)?$", orders.GetOrdersData)

	//支付
	mux.HFunc("^/l/alipay(\\?.*)?$", alipay.AlipayWebRequest)
	mux.HFilterFunc("^/api/pub/alipayWeb/notify(\\?.*)?$", alipay.AlipayWebNotify)
	mux.HFilterFunc("^/api/pub/alipayWeb/return(\\?.*)?$", alipay.AlipayWebReturn)
}



//test system is run
func test(hs *routing.HTTPSession) routing.HResult {
	type Human struct {
		Name string
		Age  int64
	}
	var (
		name string
		age int64
	)
	err := hs.ValidCheckVal(`
		name,O|S,L:0;
		age,R|I,R:0;
		`, &name, &age)
	var t = Human{}
	if nil!=C.ReErr(hs,err,C.M6){
		return routing.HRES_RETURN
	}
	err = hs.JsonObjVal2("data",&t)
	if err!=nil {
		log.E(C.M6, err)
		C.JsonRes(hs, 1, "", C.M6)
		return routing.HRES_RETURN
	}
	cookies, _ := hs.R.Cookie("token")
	fmt.Println("token:", cookies)
	log.D("all request data--:%v",t,hs.AllRVal())
	C.JsonRes(hs, 0, t, "test success")
	return routing.HRES_RETURN
}

//init uap pkg
func initUap() {
	var cfg string = ""
	cfg = "conf/uap.properties"

	err := conf.Cfg.InitWithFilePath(cfg)
	if err != nil {
		log.E("set config value err:%v",err.Error())
		return
	}
	db, err := databaseConn.GetConn()
	//initial the database
	ucs.Db = db
	if err != nil {
		log.E("connect db failde:%v",err)
		return
	}
	//initial the database
	uap.InitDb(ucs.ConDb)
	usr.CheckUcs(ucs.Db)
}

package orders
import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"encoding/json"
	D "databaseConn"
)

func AddToCartTest(hs *routing.HTTPSession) routing.HResult {
	re:=ReData{0,"",nil}
	log.D("AddToCartTest r data--:%v", hs.R.FormValue("data"))
	re = AddToCartTestF(hs,0)
	JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}


///-----------------------OLD
//获取回调地址
func GetUrl(hs *routing.HTTPSession) routing.HResult {
	re:=ReData{1,"",nil}
	db, err := D.GetConn()
	if nil!=err {
		log.E(M1, err)
		re.M = M1
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}
	log.E("传入参数", hs.R.FormValue("data"),"本地host:",hs.Host())
	data:=G_ORDERS_DATA{}
	if err := json.Unmarshal([]byte(hs.R.FormValue("data")), &data); err != nil {
		re.M = "传入数据有误！"
		log.E(re.M, err)
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}
	host:="http://"+hs.Host()
	log.E("d",data,host)
	reUrl:=""
	token,e:=GetUrlF(db,data)
	if e!=""{
		re.M = e
		JsonRes(hs, re.C, re.D, re.M)
	}
	cbApi:="/api/EditO_USR"
	reUrl = host+cbApi+"?token="+token
	re.C = 0
	re.D = reUrl
	JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}
//修改o_usr数据
func EditO_USR(hs *routing.HTTPSession) routing.HResult {
	re:=ReData{1,"",nil}
	db, err := D.GetConn()
	if nil!=err {
		log.E(M1, err)
		re.M = M1
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}

	log.E("传入参数", hs.R.FormValue("data"),hs.R.FormValue("token"))
	data:=G_ORDERS_DATA{}
	if err := json.Unmarshal([]byte(hs.R.FormValue("data")), &data); err != nil {
		re.M = "传入数据有误！"
		log.E(re.M, err)
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}
	data.Token = hs.R.FormValue("token")
	log.E("d",data)
	re.M=SyncOrdersF(db,data)
	re.C = 0
	JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}

//加入购物车或直接购买
func AddToCart(hs *routing.HTTPSession) routing.HResult {
	re:=ReData{0,"",nil}
	log.D("传入参数", hs.R.FormValue("data"))
	data:=[]*G_ADD_CART{}
	if err := json.Unmarshal([]byte(hs.R.FormValue("data")), &data); err != nil {
		re.C = 1
		re.M = "传入数据有误！"
		log.E(re.M, err)
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}

	db, err := D.GetConn()
	if nil!=err {
		log.E(M1, err)
		re.M = M1
		JsonRes(hs, re.C, re.D, re.M)
		return routing.HRES_RETURN
	}
	userId:=hs.IntVal("uid")
	re = AddToCartF(hs.R,db,data,userId)
	JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}

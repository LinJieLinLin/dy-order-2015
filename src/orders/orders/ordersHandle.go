package orders
import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/log"
	C "common"

)

//submit orders
func SubmitOrders(hs *routing.HTTPSession) routing.HResult {
	re:=C.ReData{0,"",nil}
	cartIds :=[]*C.G_SubmitData{}

	err:= hs.JsonObjVal2("data",&cartIds)
	if nil!=C.ReErr(hs,err,"AddToCart data"){
		return routing.HRES_RETURN
	}
	log.D("AddToCart data len--:%v", len(cartIds))

	re = SubmitOrderF(cartIds,hs.IntVal("uid"))

	C.JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}
// get orders data list
func GetOrdersData(hs *routing.HTTPSession) routing.HResult {
	re:=C.ReData{0,"",nil}
	data :=C.G_OrdersData{}

	err:= hs.JsonObjVal2("data",&data)
	if nil!=C.ReErr(hs,err,"AddToCart data"){
		return routing.HRES_RETURN
	}
	log.D("AddToCart data len--:%v", data)

	data.UserId = hs.IntVal("uid")
	re = GetOrdersDataF(data)

	C.JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}
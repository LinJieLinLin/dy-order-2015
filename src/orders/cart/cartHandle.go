package cart
import (
	"github.com/Centny/gwf/routing"
	C "common"
	"github.com/Centny/gwf/log"
)

//add goods to cart
func AddToCart(hs *routing.HTTPSession) routing.HResult {
	log.D("form data:--%v",hs.R.Form)
	re:=C.ReData{0,"",nil}
	data:=[]*C.O_CART{}

	err:= hs.JsonObjVal2("data",&data)
	if nil!=C.ReErr(hs,err,C.M6){
		return routing.HRES_RETURN
	}
	log.D("AddToCart input data len--:%v", len(data))

	//check User
	err=C.CheckUser(hs,data[0].USER_ID)
	if nil!=C.ReErr(hs,err,C.M2){
		return routing.HRES_RETURN
	}

	re = AddToCartF(data)
	C.JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}
//get cart data
func GetCart(hs *routing.HTTPSession) routing.HResult {
	log.D("GetCart input data--:%v")
	re:=C.ReData{0,"",nil}
	data:=C.G_CartList{}
	err:= hs.JsonObjVal2("data",&data)
	if nil!=C.ReErr(hs,err,C.M6){
		return routing.HRES_RETURN
	}
	log.D("GetCart input data--:%v", data)
	re = GetCartF(hs.IntVal("uid"),data)
	C.JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}

//EditCart data update or delete
func EditCart(hs *routing.HTTPSession) routing.HResult {
	re:=C.ReData{0,"",nil}
	data:=[]*C.G_CartData{}
	err:= hs.JsonObjVal2("data",&data)
	if nil!=C.ReErr(hs,err,C.M6){
		return routing.HRES_RETURN
	}
	log.D("EditCart input data len--:%v", len(data))
	re = EditCartF(data,hs.IntVal("uid"),nil)
	C.JsonRes(hs, re.C, re.D, re.M)
	return routing.HRES_RETURN
}

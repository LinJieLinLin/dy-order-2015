package alipay
//
//import (
//	"alipayModule"
//	"encoding/json"
//	"fmt"
//	"github.com/Centny/gwf/log"
//	"net/http"
//	"strconv"
//	"strings"
//)
//
////POST
////支付宝回调处理
//func AlipayMobileNotify(w http.ResponseWriter, r *http.Request) {
//	log.I("alipay Notify Begin")
//
//	var callbackMsg = "fail"
//	defer func() {
//		log.I("alipay Notify End")
//		log.I("callbackMsg to alipay : %v", callbackMsg)
//		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		fmt.Fprint(w, callbackMsg)
//	}()
//
//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	r.PostForm = nil
//	r.ParseForm()
//
//	log.I("==========================================================")
//	log.I("Request :%v", r)
//	log.I("==========================================================")
//
//	if err := alipayModule.VerifyMobileNotify(r, alipayModule.Alipay_mobile_config); err != nil {
//		//验证失败
//		log.E("verify notify fail")
//		return
//	}
//
//	trade_status := r.FormValue("trade_status")
//	out_trade_no := r.FormValue("out_trade_no")
//	//	buyer_email := r.FormValue("buyer_email")
//	//	subject := r.FormValue("subject")
//
//	log.I("trade_status is : %v ", trade_status)
//	log.I("out_trade_no is : %v ", out_trade_no)
//
//	var total_fee float64
//	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)
//
//	log.I("orderNos = %v ", out_trade_no)
//	if "" == out_trade_no {
//		log.E("%v", "orderNos is null")
//		return
//	}
//
//	orderNos := strings.Split(out_trade_no, ",")
//	log.D("sOrderNo_array = %v", orderNos)
//
//	if len(orderNos) == 0 {
//		log.E("%v", "orderNos size is 0")
//		return
//	}
//
//	//判断该笔订单是否在商户网站中已经做过处理
//	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
//	//如果有做过处理，不执行商户的业务程序
//
//	//注意：
//	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。
//
//	if trade_status == "TRADE_SUCCESS" {
//		if err := order_v2.PaymentForAliWeb(orderNos); nil != err {
//			log.E("PaymentForAliWap fail : %v ", err)
//			return
//		}
//		log.I("%v", "成功处理订单")
//	}
//
//	//判断是否已做操作
//
//	//判断该笔订单是否在商户网站中已经做过处理
//	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
//	//如果有做过处理，不执行商户的业务程序
//
//	//注意：
//	//1、开通了普通即时到账，买家付款成功后。
//	//该种交易状态只在两种情况下出现
//	//2、开通了高级即时到账，从该笔交易成功时间算起，过了签约时的可退款时限（如：三个月以内可退款、一年以内可退款等）后。
//
//	if trade_status == "TRADE_FINISHED" {
//
//	}
//	//	echo "success";		//请不要修改或删除
//	callbackMsg = "success"
//	return
//}
//
////获取支付宝签名
//func GetRsaSign(w http.ResponseWriter, r *http.Request) {
//
//	r.ParseForm()
//
//	log.I("orderNos : %v", r.FormValue("orderNos"))
//	log.I("price : %v", r.FormValue("price"))
//
//	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
//	if err != nil {
//		price = 0.01
//	}
//
//	amr := alipayModule.AlipayMobileRequest{}
//	amr.OutTradeNo = r.FormValue("orderNos")
//	amr.Subject = "掌上学园商品"
//	amr.Body = "掌上学园商品"
//	amr.TotalFee = price
//
//	orderinfo := alipayModule.AlipayMobileRsaSign(amr)
//
//	rs := make(map[string]interface{})
//	rs["code"] = 0
//	rs["msg"] = ""
//	rs["data"] = orderinfo //createLinkString(&p)
//
//	b, _ := json.Marshal(rs)
//	http.Error(w, string(b), 200)
//}

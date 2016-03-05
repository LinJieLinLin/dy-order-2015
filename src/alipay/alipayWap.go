package alipay
//
//import (
//	"alipayModule"
//	"config"
//	"errors"
//	"fmt"
//	"github.com/Centny/gwf/log"
//	"github.com/Centny/gwf/routing"
//	"net/http"
//	"org.cny.uas/common"
//	"strings"
//)
//
///**
//请求支付
//*/
//func AlipayWapRequest(hs *routing.HTTPSession) routing.HResult {
//
//	orderNos, err := GetOrderNosFromRequest(hs.R)
//	if err != nil {
//		return common.MsgResE(hs, 1, err.Error())
//	}
//	log.I("orderNos = %v", orderNos)
//
//	orderMsgList, err := order_v2.GetOrderMsg(orderNos)
//	if nil != err {
//		log.E("order_v2.GetOrderMsg err : %v \n", err)
//		return common.MsgResE(hs, 1, err.Error())
//	}
//
//	for _, orderMsg := range orderMsgList {
//		if orderMsg.OrderStatus != order.OrderStatusUnpaid {
//			return common.MsgResE(hs, 1, "订单已失效")
//		}
//	}
//	totalPrice, err := order_v2.GetOrdersTotalPrice(orderNos)
//	if err != nil {
//		return common.MsgResE(hs, 1, err.Error())
//	}
//	log.I("totalPrice = %v", totalPrice)
//	if float64(0) == totalPrice {
//		return common.MsgResE(hs, 1, "订单错误，请重新支付")
//	}
//
//	alipayR := alipayModule.AlipayWebRequest{
//		OutTradeNo: strings.Join(orderNos, ","), // 订单号
//		Subject:    `云平台资源`, // 商品名称
//		TotalFee:   0.01, // 价格
//	}
//
//	// 输出的是 html 页面，会自动跳转到支付界面
//	err = alipayModule.AlipayWapRequest(*alipayModule.Alipay_wap_config, alipayR, hs.W)
//	if err != nil {
//		return common.MsgResE(hs, 1, err.Error())
//	}
//	return common.MsgRes(hs, "")
//}
//
////支付宝异步通知处理
//func AlipayWapNotify(w http.ResponseWriter, r *http.Request) {
//	log.I("AlipayWapNotify Begin")
//
//	var callbackMsg = "fail"
//	defer func() {
//		log.I("AlipayWapNotify Notify End")
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
//	log.I("AlipayWapNotify Request :%v", r)
//	log.I("==========================================================")
//
//	alipayNotify, err := alipayModule.VerifyWapNotify(r, alipayModule.Alipay_wap_config)
//	if err != nil {
//		//验证失败
//		log.E("verify notify fail")
//		return
//	}
//
//	trade_status := alipayNotify.Trade_status
//	out_trade_no := alipayNotify.Out_trade_no
//	buyer_email := alipayNotify.Buyer_email
//	subject := alipayNotify.Subject
//
//	log.I("trade_status is : %v ", trade_status)
//	log.I("out_trade_no is : %v ", out_trade_no)
//	log.I("buyer_email is : %v ", buyer_email)
//	log.I("subject is : %v ", subject)
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
//
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
//		if err := order_v2.PaymentForAliWeb(orderNos); nil != err {
//			log.E("PaymentForAliWap fail : %v ", err)
//			return
//		}
//		log.I("%v", "成功处理订单")
//	}
//	//	echo "success";		//请不要修改或删除
//	callbackMsg = "success"
//	return
//}
//
////支付宝 同步通知处理
//func AlipayWapCallback(w http.ResponseWriter, r *http.Request) {
//	log.I("AlipayWapCallback Begin")
//
//	var callbackMsg = "fail"
//	defer func() {
//		log.I("AlipayWapCallback End")
//		log.I("callbackMsg to alipay : %v", callbackMsg)
//		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		fmt.Fprint(w, callbackMsg)
//	}()
//
//	//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	//	r.PostForm = nil
//	r.ParseForm()
//
//	log.I("==========================================================")
//	log.I("AlipayWapCallback Request :%v", r)
//	log.I("==========================================================")
//
//	if err := alipayModule.VerifyWapCallback(r, alipayModule.Alipay_wap_config); err != nil {
//		//验证失败
//		log.E("verify notify fail")
//		return
//	}
//
//	trade_status := r.FormValue("result")
//	out_trade_no := r.FormValue("out_trade_no")
//
//	log.I("trade_status is : %v ", trade_status)
//	log.I("out_trade_no is : %v ", out_trade_no)
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
//	var total_fee float64
//	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)
//
//	if trade_status == "success" {
//
//		orderMsgList, err := order_v2.GetOrderMsg(orderNos)
//		if nil != err {
//			log.E("payment.GetOrderMsg err : %v \n", err)
//			return
//		}
//
//		var rechargeList []string = []string{}
//		var shoppingList []string = []string{}
//
//		for _, orderMsg := range orderMsgList {
//			if orderMsg.OrderType == "SHOPPING" {
//				shoppingList = append(shoppingList, orderMsg.OrderNo)
//			}
//			if orderMsg.OrderType == "RECHARGE" {
//				rechargeList = append(rechargeList, orderMsg.OrderNo)
//			}
//		}
//
//		if len(shoppingList) > 0 {
//			http.Redirect(w, r, config.AlipayNotifyHost() + alipayModule.Alipay_wap_config.Show_order_url + "?type=shopping&orderNos=" + out_trade_no, http.StatusFound)
//		}
//
//		if len(rechargeList) > 0 {
//			http.Redirect(w, r, config.AlipayNotifyHost() + alipayModule.Alipay_wap_config.Show_order_url + "?type=recharge&orderNos=" + out_trade_no, http.StatusFound)
//		}
//	}
//	callbackMsg = "success"
//	return
//}
//
//func GetOrderNosFromRequest(r *http.Request) (orderNos []string, err error) {
//
//	r.ParseForm()
//
//	tmpOrderNo := r.FormValue("orderNos")
//	log.I("orderNos = %v ", tmpOrderNo)
//	if "" == tmpOrderNo {
//		log.E("%v", "orderNos is null")
//		err = errors.New("orderNos is null")
//		return
//	}
//
//	orderNos = strings.Split(tmpOrderNo, ",")
//	log.D("sOrderNo_array = %v", orderNos)
//
//	if len(orderNos) == 0 {
//		err = errors.New("orderNos is null ")
//		return
//	}
//
//	return
//}

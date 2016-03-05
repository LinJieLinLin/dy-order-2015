package alipay

import (
	"alipayModule"
//	"config"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
//	"net/http"
	"org.cny.uas/common"
	C "common"
	"fmt"
	"config"
	"net/http"
)

/**
请求支付
*/
func AlipayWebRequest(hs *routing.HTTPSession) routing.HResult {
	re := C.ReData{1, "", nil}
	log.D("传入参数", hs.R.FormValue("publicOrderNo"))
	orderNo := hs.R.FormValue("orderNo");
	if orderNo=="" {
		re.M="参数为空"
		C.JsonRes(hs, re.C, re.D, re.M)
	}
	//检查订单信息返回资源价格
	price, e := GetOrdersPrice(orderNo)
	log.E("123",orderNo,e)
	if e!="" {
		w:=hs.W
		r:=hs.R
		http.Redirect(w, r, config.AlipayNotifyHost()+alipayModule.Alipay_web_config.Show_order_url+"?type=fail&msg="+e+"&orderNos="+orderNo, http.StatusFound)
		return routing.HRES_RETURN
	}
	price = 0.01
	alipayR := alipayModule.AlipayWebRequest{
		OutTradeNo: orderNo, // 订单号
		Subject:    `大洋资源`, // 商品名称
		TotalFee:   price, // 价格
	}
	//
	//	// 输出的是 html 页面，会自动跳转到支付界面
	err := alipayModule.AlipayWebRequestForm(*alipayModule.Alipay_web_config, alipayR, hs.W)
	if err != nil {
		return common.MsgResE(hs, 1, err.Error())
	}
	return routing.HRES_RETURN
}

////支付宝异步通知处理
func AlipayWebNotify(hs *routing.HTTPSession) routing.HResult {
	log.I("AlipayWebNotify Begin")
	log.E("------------------\n用户ID：", hs.StrVal("uid"))
	w := hs.W
	r := hs.R
	var callbackMsg = "fail"
	defer func() {
		log.I("AlipayWebNotify Notify End")
		log.I("callbackMsg to alipay : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	log.I("==========================================================")
	log.I("AlipayWebNotify Request :%v", r)
	log.I("==========================================================")

	if err := alipayModule.VerifyWebNotify(r, alipayModule.Alipay_web_config); err != nil {
		//验证失败
		log.E("verify notify fail")
		return routing.HRES_RETURN
	}
	trade_status := r.FormValue("trade_status")
	out_trade_no := r.FormValue("out_trade_no")
	buyer_email := r.FormValue("buyer_email")
	subject := r.FormValue("subject")

	log.I("trade_status is : %v ", trade_status)
	log.I("out_trade_no is : %v ", out_trade_no)
	log.I("buyer_email is : %v ", buyer_email)
	log.I("subject is : %v ", subject)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	log.I("orderNos = %v ", out_trade_no)
	if "" == out_trade_no {
		log.E("%v", "orderNos is null")
		return routing.HRES_RETURN
	}

	orderNos := out_trade_no
	log.D("sOrderNo_array = %v", orderNos)

	if len(orderNos) == 0 {
		log.E("%v", "orderNos size is 0")
		return routing.HRES_RETURN
	}

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。

	if trade_status == "TRADE_SUCCESS" {
		//这里进行订单处理：
		if err := EditOrders(orderNos); err != "" {
			log.E("AlipayWebNotify fail : %v ", err)
			if err!="订单已支付！" {
				return routing.HRES_RETURN
			}
		}
		log.I("%v", "支付宝异步通知处理成功处理订单")
	}

	//判断是否已做操作

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//1、开通了普通即时到账，买家付款成功后。
	//该种交易状态只在两种情况下出现
	//2、开通了高级即时到账，从该笔交易成功时间算起，过了签约时的可退款时限（如：三个月以内可退款、一年以内可退款等）后。

	if trade_status == "TRADE_FINISHED" {

	}
	//	echo "success";		//请不要修改或删除
	callbackMsg = "success"
	return routing.HRES_RETURN
}

//支付宝 同步通知处理
func AlipayWebReturn(hs *routing.HTTPSession) routing.HResult {
	log.I("AlipayWebReturn Begin")
	log.E("------------------\n用户ID：", hs.StrVal("uid"))
	w := hs.W
	r := hs.R
	var callbackMsg = "fail"
	defer func() {
		log.I("AlipayWebReturn End")
		log.I("callbackMsg to alipay : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	r.PostForm = nil
	r.ParseForm()

	log.I("==========================================================")
	log.I("AlipayWebReturn Request :%v", r)
	log.I("==========================================================")

	if err := alipayModule.VerifyWebNotify(r, alipayModule.Alipay_web_config); err != nil {
		//验证失败
		log.E("verify notify fail")
		callbackMsg = "verify notify fail"
		callbackMsg = "网页已过期"
		return routing.HRES_RETURN
	}

	trade_status := r.FormValue("trade_status")
	out_trade_no := r.FormValue("out_trade_no")
	buyer_email := r.FormValue("buyer_email")
	subject := r.FormValue("subject")
	log.I("buyer_email is : %v ", buyer_email)
	log.I("subject is : %v ", subject)
	log.I("trade_status is : %v ", trade_status)
	log.I("out_trade_no is : %v ", out_trade_no)
	if "" == out_trade_no {
		log.E("%v", "orderNos is null")
		callbackMsg = "orderNos is null"
		return routing.HRES_RETURN
	}

	orderNos := out_trade_no
	log.D("sOrderNo_array = %v", orderNos)

	if len(orderNos) == 0 {
		log.E("%v", "orderNos size is 0")
		callbackMsg = "orderNos size is 0"
		return routing.HRES_RETURN
	}

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	if trade_status == "TRADE_SUCCESS" {
		//这里进行订单处理：
		if err := EditOrders(orderNos); err != "" {
			log.E("AlipayWebNotify fail : %v ", err)
			if err!="订单已支付！" {
				http.Redirect(w, r, config.AlipayNotifyHost()+alipayModule.Alipay_web_config.Show_order_url+
				"?type=fail&orderNos="+out_trade_no+"&msg="+err, http.StatusFound)
				return routing.HRES_RETURN
			}
		}
		log.I("%v", "同步通知处理成功处理订单")
		//返回订单系统页面
		http.Redirect(w, r, config.AlipayNotifyHost()+alipayModule.Alipay_web_config.Show_order_url+"?type=success&orderNos="+out_trade_no, http.StatusFound)
		//
	}

	if trade_status == "TRADE_FINISHED" {

	}
	callbackMsg = "success"
	return routing.HRES_RETURN
}

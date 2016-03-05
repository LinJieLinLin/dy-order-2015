package alipayModule

import (
	"config"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var alipayGatewayNew = `https://mapi.alipay.com/gateway.do?`

/**
 * WAP形式消息验证地址
 */
var wap_alipayGatewayNew = `http://wappaygw.alipay.com/service/rest.htm?`

/**
 * HTTPS形式消息验证地址
 */
var https_verify_url = "https://mapi.alipay.com/gateway.do?service=notify_verify&"

/**
 * HTTP形式消息验证地址
 */
var http_verify_url = "http://notify.alipay.com/trade/notify_query.do?"

//for the test
var AlipayVerify = true

func AlipayWebRequestForm(alipay_config Alipay_config_struct, r AlipayWebRequest, w io.Writer) error {
	p := Kvpairs{
		Kvpair{`total_fee`, fmt.Sprintf("%.2f", r.TotalFee)},
		Kvpair{`subject`, r.Subject},
		//		Kvpair{`body`, r.Body},
		//		Kvpair{`show_url`, r.ShowUrl},
		Kvpair{`out_trade_no`, r.OutTradeNo},
		Kvpair{`service`, alipay_config.Service},
		Kvpair{`partner`, alipay_config.Partner},
		Kvpair{`payment_type`, alipay_config.Payment_type},
		Kvpair{`notify_url`, config.AlipayNotifyHost() + alipay_config.Notify_url},
		Kvpair{`return_url`, config.AlipayNotifyHost() + alipay_config.Return_url},
		Kvpair{`seller_email`, alipay_config.Seller_id},
		Kvpair{`_input_charset`, alipay_config.Input_charset},
	}
	paraFilter(&p)

	argSort(&p)
	sign := md5Sign(createLinkStringNoUrl(&p), alipay_config.Key)

	p = append(p, Kvpair{`sign`, sign})
	p = append(p, Kvpair{`sign_type`, `MD5`})

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='post'> `, alipayGatewayNew)
	for _, kv := range p {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	}
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
	return nil
}

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWapNotify(r *http.Request, alipay_config *Alipay_config_struct) (Notify, error) {
	log.I("VerifyWapNotify begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	notify_data := ""
	var alipayNotify Notify

	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		case "notify_data":
			notify_data = v
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	if notify_data == "" {
		log.I("notify_data is null")
		return alipayNotify, errors.New("notify_data is null")
	}

	err := xml.Unmarshal([]byte(notify_data), &alipayNotify)
	if err != nil {
		log.E("xml Unmarshal err : %v", err)
		return alipayNotify, err
	}

	log.I("p = %v", p)
	p2 := sortNotifyPara(p)
	log.I("p2 = %v", p2)
	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(&p2)

	log.I("VerifyWapNotify prestr is : %v ", prestr)
	log.I("VerifyWapNotify sign is : %v  , sign_type is %v", sign, sign_type)

	switch alipay_config.Sign_type {
	case "MD5":
		if md5Sign(prestr, alipay_config.Key) != sign {
			return alipayNotify, fmt.Errorf("sign invalid")
		}
		break
	default:
		return alipayNotify, fmt.Errorf("no right sign_type")
	}

	log.I("VerifyWapNotify Notify is %v", alipayNotify)

	notify_id := alipayNotify.Notify_id
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipay_config)
	if err != nil {
		return alipayNotify, err
	}
	log.I("VerifyWapNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.E("VerifyWapNotify responseTxt verify fail ")
		return alipayNotify, fmt.Errorf("VerifyWapNotify responseTxt is wrong")
	}
	log.I("VerifyWapNotify responseTxt verify success ")
	return alipayNotify, nil
}

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWapCallback(r *http.Request, alipay_config *Alipay_config_struct) error {
	log.I("VerifyWapCallback begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	log.I("VerifyWapCallback p = %v", p)
	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.I("VerifyWapCallback prestr is : %v ", prestr)
	log.I("VerifyWapCallback sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipay_config.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.I("VerifyWapCallback success")
	return nil
}

/**
 * 异步通知时，对参数做固定排序
 * @param $para 排序前的参数组
 * @return 排序后的参数组
 */
func sortNotifyPara(para *Kvpairs) Kvpairs {
	new := Kvpairs{}
	for _, kv := range *para {
		if kv.K == "service" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "v" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "sec_id" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "notify_data" {
			new = append(new, kv)
		}
	}
	return new
}

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWebReturn(r *http.Request, config *Alipay_config_struct) error {
	log.I("VerifyWebReturn begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.PostForm.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.I("VerifyWebReturn prestr is : %v ", prestr)
	log.I("VerifyWebReturn sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, config.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.I("VerifyWebReturn success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, config)
	if err != nil {
		return err
	}
	log.I("VerifyWebReturn responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.E("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.I("VerifyWebReturn responseTxt verify success ")
	return nil

}

/**
 * 针对notify_url验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyMobileNotify(r *http.Request, config *Alipay_config_struct) error {

	log.I("VerifyMobileNotify begin")

	signErr := verifySign(r.PostForm, config)
	if signErr != nil {
		return signErr
	}
	log.I("VerifyMobileNotify verifySign success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）
	responseTxt, err := getResponse(notify_id, config)
	if err != nil {
		return err
	}
	log.I("VerifyMobileNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.I("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.I("VerifyMobileNotify responseTxt verify success ")
	return nil
}

/**
 * 针对notify_url验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWebNotify(r *http.Request, config *Alipay_config_struct) error {
	log.I("VerifyWebNotify begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.I("VerifyWebNotify prestr is : %v ", prestr)
	log.I("VerifyWebNotify sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, config.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.I("VerifyWebNotify success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, config)
	if err != nil {
		return err
	}
	log.I("VerifyWebNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.E("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.I("VerifyWebNotify responseTxt verify success ")
	return nil
}

func RsaSign(para *Kvpairs, config *Alipay_config_struct) string {
	buildRequestPara(para, config)
	// buildRequestMysign(para, config)
	return createLinkstringUrlencode(para)
}

func AlipayMobileRsaSign(amr AlipayMobileRequest) string {

	p := Kvpairs{
		Kvpair{`_input_charset`, Alipay_mobile_config.Input_charset},
		Kvpair{`partner`, Alipay_mobile_config.Partner},
		Kvpair{`payment_type`, Alipay_mobile_config.Payment_type},
		Kvpair{`notify_url`, config.AlipayNotifyHost() + Alipay_mobile_config.Notify_url},
		Kvpair{`service`, Alipay_mobile_config.Service},
		Kvpair{`seller_id`, Alipay_mobile_config.Seller_id},
		Kvpair{`out_trade_no`, amr.OutTradeNo},
		Kvpair{`subject`, amr.Subject},
		Kvpair{`total_fee`, fmt.Sprintf("%.2f", amr.TotalFee)},
		Kvpair{`body`, amr.Body},
		//				Kvpair{`it_b_pay`, `15m`},
	}
	RsaSign(&p, Alipay_mobile_config)
	return createLinkstringUrlencode(&p)
}

func AlipayWapRequest(alipay_config Alipay_config_struct, r AlipayWebRequest, w io.Writer) error {

	var format = "xml"
	var v = "2.0"
	args := []interface{}{}
	args = append(args, config.AlipayNotifyHost()+alipay_config.Notify_url)
	//	args = append(args,url.QueryEscape(alipay_config.Notify_url) )
	args = append(args, config.AlipayNotifyHost()+alipay_config.Wap_callback_url)
	//	args = append(args, url.QueryEscape(alipay_config.Wap_callback_url))
	args = append(args, alipay_config.Seller_id)
	args = append(args, r.OutTradeNo)
	args = append(args, r.Subject)
	args = append(args, fmt.Sprintf("%.2f", r.TotalFee))
	args = append(args, config.AlipayNotifyHost()+alipay_config.Wap_merchant_url)
	//	args = append(args, url.QueryEscape(alipay_config.Wap_merchant_url))

	var token_req_data = fmt.Sprintf("<direct_trade_create_req><notify_url>%v</notify_url><call_back_url>%v</call_back_url><seller_account_name>%v</seller_account_name><out_trade_no>%v</out_trade_no><subject>%v</subject><total_fee>%v</total_fee><merchant_url>%v</merchant_url></direct_trade_create_req>", args...)

	log.I("token_req_data = %v", token_req_data)
	p := Kvpairs{
		Kvpair{`service`, alipay_config.Wap_Service},
		Kvpair{`partner`, alipay_config.Partner},
		Kvpair{`sec_id`, alipay_config.Sign_type},
		Kvpair{`format`, format},
		Kvpair{`v`, v},
		Kvpair{`req_id`, r.OutTradeNo},
		Kvpair{`req_data`, token_req_data},
		Kvpair{`_input_charset`, alipay_config.Input_charset},
	}
	paraFilter(&p)

	argSort(&p)

	buildRequestPara(&p, &alipay_config)

	log.I("pararms is  = %v", p)

	origin_token, err := getHttpResponsePOST(wap_alipayGatewayNew, "", alipay_config.Input_charset, &p)
	if nil != err {
		log.E("get token fail : %v", err)
		return err
	}
	//	fmt.Fprintln(w, origin_token)
	//	return nil
	log.I("origin_token = %v", origin_token)
	origin_token, _ = url.QueryUnescape(origin_token)
	log.I("url.QueryUnescape origin_token = %v", origin_token)

	token := ParseOriginTokenMsg(origin_token, &alipay_config)
	if token == "" {
		return errors.New("解析token失败")
	}

	var trade_req_data = fmt.Sprintf("<auth_and_execute_req><request_token>%v</request_token></auth_and_execute_req>", token)

	p2 := Kvpairs{
		Kvpair{`service`, alipay_config.Service},
		Kvpair{`partner`, alipay_config.Partner},
		Kvpair{`sec_id`, alipay_config.Sign_type},
		Kvpair{`format`, format},
		Kvpair{`v`, v},
		Kvpair{`req_id`, r.OutTradeNo},
		Kvpair{`req_data`, trade_req_data},
		Kvpair{`_input_charset`, alipay_config.Input_charset},
	}
	paraFilter(&p2)

	argSort(&p2)

	buildRequestPara(&p2, &alipay_config)

	log.I("pararm2 = %v", p2)

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='get'> `, wap_alipayGatewayNew)
	for _, kv := range p2 {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	}
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
	return nil
}

func ParseOriginTokenMsg(origin_token string, config *Alipay_config_struct) string {

	tokenArray := strings.Split(origin_token, "&")
	kvs := Kvpairs{}
	tokenkv := Kvpair{}
	for _, v := range tokenArray {
		log.I("v : %v", v)
		lenght := len(v)
		index := strings.Index(v, "=")
		if index != -1 {
			str1 := Substr(v, index+1, lenght-index-1)
			str2 := Substr(v, 0, index)
			log.I("str1 = %v", str1)
			log.I("str2 = %v", str2)
			subkv := Kvpair{str2, str1}
			kvs = append(kvs, subkv)
			if str2 == "res_data" {
				tokenkv = Kvpair{str2, str1}
			}
		}
	}

	//	token_data_Decrypt := ""

	//	for _, v := range kvs {
	//		if v.K == "sign_type" && v.V == "MD5" {
	//			token_data_Decrypt, err := rsaDecrypt(tokenkv.V, config.Private_key)
	//			if nil != err {
	//				return token_data_Decrypt
	//			}
	//		}
	//	}

	log.I("kvs = %v", kvs)

	log.I("tokenkv.V = %v", tokenkv.V)

	if tokenkv.V == "" {
		return ""
	}

	type Direct_trade_create_res struct {
		Token string `xml:"request_token"`
	}

	var alipay Direct_trade_create_res

	err := xml.Unmarshal([]byte(tokenkv.V), &alipay)
	if err != nil {
		log.E("xml Unmarshal err : %v", err)
		return ""
	}

	log.I("alipay xml : %v", alipay)

	if alipay.Token == "" {
		log.E("token is null  ")
		return ""
	}

	return alipay.Token
}

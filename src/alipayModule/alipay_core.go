/**
 * Created with IntelliJ IDEA.
 * User: liangjunyu
 * Date: 14-2-18
 * Time: 下午2:11
 * To change this template use File | Settings | File Templates.
 */
package alipayModule

import (
	"crypto/tls"
	"fmt"
	"github.com/Centny/gwf/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/**
 * 远程获取数据，POST模式
 * 注意：
 * 1.使用Crul需要修改服务器中php.ini文件的设置，找到php_curl.dll去掉前面的";"就行了
 * 2.文件夹中cacert.pem是SSL证书请保证其路径有效，目前默认路径是：getcwd().'\\cacert.pem'
 * @param $url 指定URL完整路径地址
 * @param $cacert_url 指定当前工作目录绝对路径
 * @param $para 请求的数据
 * @param $input_charset 编码格式。默认值：空值
 * return 远程输出的数据
 */
func getHttpResponsePOST(surl, cacert_url, charSet string, kv *Kvpairs) (string, error) {
	//todo :  getHttpResponseGET

	//	tr := &http.Transport{
	//		TLSClientConfig:
	//		&tls.Config{InsecureSkipVerify: true},
	//	}
	values := createLinkstringForPost(kv)

	//	surl = surl + "_input_charset=" + charSet

	log.I("strings.NewReader(values) :%v", strings.NewReader(values))
	req, _ := http.NewRequest("GET", surl + values, strings.NewReader("")) //, strings.NewReader(values))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//	client := &http.Client{Transport: tr}
	//	client := &http.Client{}
	resp, err := http.DefaultClient.Do(req)

	//	values2:= url.Values{}
	//	for _, kv := range *kv {
	//		if kv.K == "notify_url" || kv.K == "sign" {
	//			//移动支付这里要做URL转换，加双引号  坑
	////			strs = append(strs, kv.K + "=" + url.QueryEscape(kv.V) + "")
	//			values2.Add(kv.K,url.QueryEscape(kv.V))
	//		} else {
	////			strs = append(strs, kv.K + "=" + kv.V + "")
	//			values2.Add(kv.K,kv.V)
	//		}
	//	}

	//	resp, err := http.DefaultClient.PostForm(surl,values2)
	if err != nil {
		log.I("http.Client get err : %v", err)
		return "", err
	}
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	return string(bodyByte), nil
}

/**
 * 远程获取数据，GET模式
 * 注意：
 * 1.使用Crul需要修改服务器中php.ini文件的设置，找到php_curl.dll去掉前面的";"就行了
 * 2.文件夹中cacert.pem是SSL证书请保证其路径有效，目前默认路径是：getcwd().'\\cacert.pem'
 * @param $url 指定URL完整路径地址
 * @param $cacert_url 指定当前工作目录绝对路径
 * return 远程输出的数据
 */
func getHttpResponseGET(url, cacert_url string) (string, error) {
	log.I("http.Client get url : %v", url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.E("http.Client get err : %v", err)
		return "", err
	}
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	log.I("http.Client get bodyByte : %v", string(bodyByte))
	if !AlipayVerify {
		return `true`, nil
	}
	return string(bodyByte), nil
}

/**
 * 实现多种字符编码方式
 * @param $input 需要编码的字符串
 * @param $_output_charset 输出的编码格式
 * @param $_input_charset 输入的编码格式
 * return 编码后的字符串
 */
func charsetEncode() {}

/**
 * 实现多种字符解码方式
 * @param $input 需要解码的字符串
 * @param $_output_charset 输出的解码格式
 * @param $_input_charset 输入的解码格式
 * return 解码后的字符串
 */
func charsetDecode() {}

/**
 * 生成签名结果
 * @param $para_sort 已排序要签名的数组
 * return 签名结果字符串
 */
func buildRequestMysign(kv *Kvpairs, config *Alipay_config_struct) string {
	mysign := ""
	var err error
	switch config.Sign_type {
	case "RSA":
		prestr := createLinkstringUrlencode(kv)
		mysign, err = rsaSign(prestr, config.Private_key)
		if err != nil {
			mysign = "rsaSign err"
		}
		break
	case "0001":
		prestr := createLinkStringNoUrl(kv)
		mysign, err = rsaSign(prestr, config.Private_key)
		if err != nil {
			mysign = "rsaSign err"
		}
		break
	case "MD5":
		return md5Sign(createLinkStringNoUrl(kv), config.Key)
		break
	default:
		mysign = ""
	}
	return mysign
}

/**
 * 生成要请求给支付宝的参数数组
 * @param $para_temp 请求前的参数数组
 * @return 要请求的参数数组
 */
func buildRequestPara(kv *Kvpairs, config *Alipay_config_struct) {
	//除去待签名参数数组中的空值和签名参数
	paraFilter(kv)

	//对待签名参数数组排序
	argSort(kv)

	//生成签名结果
	mysign := buildRequestMysign(kv, config)

	//签名结果与签名方式加入请求提交参数组中
	*kv = append(*kv, Kvpair{`sign`, mysign})

	for _, kvSingle := range *kv {
		if kvSingle.K == "service" {
			if kvSingle.V != "alipay.wap.trade.create.direct" && kvSingle.V != "alipay.wap.auth.authAndExecute" {
				*kv = append(*kv, Kvpair{`sign_type`, config.Sign_type})
				break
			}
		}
	}
}

/**
 * 获取返回时的签名验证结果
 * @param $para_temp 通知返回来的参数数组
 * @param $sign 返回的签名结果
 * @return 签名验证结果
 */
func verifySign(u url.Values, config *Alipay_config_struct) error {

	log.I("verify sign values :%v", u)
	p := &Kvpairs{}
	sign := ""
	for k := range u {
		v := u.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
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

	encodeString, err := base64EnCode(sign)
	if err != nil {
		log.E("base64EnCode err : %v", err)
	}
	log.I("prestr is : %v ", prestr)
	log.I("sign is : %v  , sign_type is %v", encodeString, config.Sign_type)

	switch config.Sign_type {
	case "RSA":
		return rsaVerify(prestr, encodeString, config.Public_key)
	}
	return fmt.Errorf("not right config.Sign_type ： %v", config.Sign_type)
}

/**
 * 获取远程服务器ATN结果,验证返回URL
 * @param $notify_id 通知校验ID
 * @return 服务器ATN结果
 * 验证结果集：
 * invalid命令参数不对 出现这个错误，请检测返回处理中partner和key是否为空
 * true 返回正确信息
 * false 请检查防火墙或者是服务器阻止端口问题以及验证时间是否超过一分钟
 */
func getResponse(notify_id string, config *Alipay_config_struct) (string, error) {
	transport := config.Transport
	partner := config.Partner
	verify_url := ""
	if transport == "https" {
		verify_url = https_verify_url
	} else {
		verify_url = http_verify_url
	}

	verify_url = verify_url + "partner=" + partner + "&notify_id=" + notify_id
	return getHttpResponseGET(verify_url, config.Cacert)
}

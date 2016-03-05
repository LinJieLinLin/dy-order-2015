package alipay
//
//import (
//	//	"databaseConn"
//	"alipayModule"
//	. "github.com/smartystreets/goconvey/convey"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//	"strings"
//	"testing"
//)
//
//// go test -coverprofile=c.out
//// go tool cover -html=c.out -o c.html
//
//var db_string = "cny:123@tcp(192.168.1.17:3306)/rcp_test_01?charset=utf8"
//
//func TestAlipayNotify(t *testing.T) {
//
//	alipayModule.Init(alipayModule.Alipay_wap_config)
//
//	Vstring := ""
//	values := url.Values{
//		"is_total_fee_adjust": {"N"},
//		"use_coupon":          {"N"},
//		"body":                {"教育云平台资源"},
//		"subject":             {"教育云平台资源"},
//		"price":               {"0.01"},
//		"notify_id":           {"cc8d1959a88792ffbc9d4994c6b5f29a5k"},
//		"payment_type":        {"1"},
//		"quantity":            {"1"},
//		"trade_status":        {"TRADE_FINISHED"},
//		"sign_type":           {"RSA"},
//		"notify_time":         {"2014-02-27 16:10:40"},
//		"seller_email":        {"itdayang@gmail.com"},
//		"gmt_close":           {"2014-02-27 16:10:40"},
//		"gmt_create":          {"2014-02-27 16:10:40"},
//		"gmt_payment":         {"2014-02-27 16:10:40"},
//		"out_trade_no":        {"201402271610290"},
//		"buyer_id":            {"2088302084610640"},
//		"buyer_email":         {"619167142@qq.com"},
//		"total_fee":           {"0.01"},
//		"trade_no":            {"2014022717326064"},
//		"notify_type":         {"trade_status_sync"},
//		"discount":            {"0.00"},
//		"seller_id":           {"2088501949844011"},
//		"sign":                {url.QueryEscape("eTcAzyHDZdX907cUub1oGJYbtfu9zayOPju5R9RKAmETtp56DLalvVQBPzj/6oIlZC+wUtzTIqVSP4FsD8PYyR/Tc9BcrJTC8/wCzpE2J1bub/9fcaggWBFMGoFZxgGcWu+0l2BJat0xUIAclL6KpItd4S6+LM9xmJAg9lN226c=")},
//	}
//
//	for i, v := range values {
//		tmps := i
//		Vstring += tmps
//		Vstring += "="
//		tmps = v[0]
//		Vstring += tmps
//		Vstring += "&"
//	}
//
//	Convey("deal alipay verify success", t, func() {
//			defer func() {
//				alipayModule.AlipayVerify = true
//			}()
//			alipayModule.AlipayVerify = false
//
//			handler := http.HandlerFunc(Notify)
//			req, _ := http.NewRequest("POST", "", strings.NewReader(Vstring))
//			resqw := httptest.NewRecorder()
//
//			handler.ServeHTTP(resqw, req)
//			resp_body := resqw.Body.String()
//
//			So(resp_body, ShouldEqual, "fail")
//		})
//
//	Convey("deal alipay get the alipay resp fail", t, func() {
//
//			handler := http.HandlerFunc(Notify)
//			req, _ := http.NewRequest("POST", "", strings.NewReader(Vstring))
//			resqw := httptest.NewRecorder()
//
//			handler.ServeHTTP(resqw, req)
//			resp_body := resqw.Body.String()
//
//			So(resp_body, ShouldEqual, "fail")
//		})
//
//	Convey("deal alipay not verify success", t, func() {
//			Vstring_backup := Vstring + "ddd"
//			handler := http.HandlerFunc(Notify)
//			req, _ := http.NewRequest("POST", "", strings.NewReader(Vstring_backup))
//			resqw := httptest.NewRecorder()
//
//			handler.ServeHTTP(resqw, req)
//			resp_body := resqw.Body.String()
//
//			So(resp_body, ShouldEqual, "fail")
//		})
//
//}

package common
import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/log"
	"databaseConn"
	"github.com/Centny/gwf/util"
	"encoding/json"
	"database/sql"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"fmt"
	"github.com/Centny/gwf/dbutil"
	"strconv"
)
// err return when hs is exist
func ReErr(hs *routing.HTTPSession, arg_err error, arg_msg string) error {
	if nil!=arg_err {
		log.E(arg_msg+"--:%v", arg_err)
		JsonRes(hs, 1, "", arg_msg)
	}
	return arg_err
}
//err return for ***func
func ReFuncErr(re *ReData, arg_err error, arg_msg string) error {
	if nil!=arg_err {
		log.E(arg_msg+"--:%v", arg_err)
		re.M = arg_msg
	}
	return arg_err
}
//err return for ***func with tx
func ReFuncErrTx(tx *sql.Tx, re *ReData, arg_err error, arg_msg string) error {
	if nil!=arg_err {
		log.E(arg_msg+"--:%v", arg_err)
		re.M = arg_msg
	}
	return arg_err
}
//check orders table if not exist create
func CheckT(db *sql.DB) error {
	if nil==db {
		return errors.New("db conn fail")
	}
	if _, err := dbutil.DbQuery(db, selectOCart, 1); err != nil {
		log.D("orders table not found, auto creating...")
		err = dbutil.DbExecScript(db, ordersTable)
		if err != nil {
			log.I("creat orders table failed--:%v", err)
			return err
		}
	}
	return nil
}
//set Session
func SessionFilter(hs *routing.HTTPSession) routing.HResult {
	log.D("login userId is--:%v", hs.StrVal("uid"))
	if "" == hs.StrVal("uid") {
		return routing.HRES_RETURN
	}
	if ""== hs.StrVal("username") {
		db, err := databaseConn.GetConn()
		if nil!=ReErr(hs, err, M1) {
			return routing.HRES_RETURN
		}
		sql := selectUserName
		userName := ""
		log.D("get username sql--:%v--:%v", sql, hs.StrVal("uid"))
		if err := db.QueryRow(sql, hs.StrVal("uid")).Scan(&userName); err != nil {
			log.E(M2, err)
			JsonRes(hs, 0, nil, M2)
			return routing.HRES_RETURN
		}
		hs.SetVal("username", userName)
		log.D("login userName is--:%v", userName)
	}
	return routing.HRES_CONTINUE
}
//quit system
func Quit(hs *routing.HTTPSession) routing.HResult {
	hs.SetVal("uid", nil)
	hs.SetVal("username", nil)
	log.D("now quit system")
	return routing.HRES_CONTINUE
}
//return json data
func JsonRes(hs *routing.HTTPSession, code int64, data interface{}, msg string) routing.HResult {
	res := make(util.Map)
	res["code"] = code
	res["msg"] = msg
	res["data"] = data
	dBys, _ := json.Marshal(res)
	hs.W.Write(dBys)
	return routing.HRES_RETURN
}
//check user is login
func CheckUser(hs *routing.HTTPSession, arg_userId int64) error {
	err := ""
	if "" == hs.StrVal("uid") {
		err="用户未登陆！"
	}
	if hs.IntVal("uid")!=arg_userId {
		err="用户ID不匹配，请重新登陆！"
	}
	if ""!=err {
		return errors.New(err)
	}
	return nil
}
//mock test request with hs
func TestRequestWithHs(arg_mux *routing.SessionMux, arg_method string, arg_url string, arg_v url.Values, cookie []*http.Cookie) (ReData) {
	fmt.Println("URL--:%v--%v--%v", arg_url, arg_v,arg_method)
	v:=arg_v.Encode()
	if "GET"== arg_method||"get"==arg_method{
		arg_url+="?"+v
	}
	r, _ := http.NewRequest(arg_method, arg_url, ioutil.NopCloser(strings.NewReader(v)))
	if "POST" == arg_method {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}
	for _, v := range cookie {
		r.AddCookie(v)
	}

//	fmt.Println("请求头：",r)

	w := httptest.NewRecorder()
	arg_mux.ServeHTTP(w, r)

	re := ReData{}
	resp_body := w.Body.String()
	fmt.Println("返回数据-----：%v",resp_body)
	err := util.Json2S(resp_body, &re)
	if nil!=err {
		fmt.Println("json to struct failed--:%v", err)
		re.M = resp_body
	}
	return re
}
//mock test request with w r
func TestRequestWithRW(fun http.HandlerFunc, arg_method string, arg_url string, arg_v url.Values, cookie []*http.Cookie) (ReData) {
	handler := http.HandlerFunc(fun)
	fmt.Println("URL--:%v--%v--%v", arg_url, arg_v,arg_method)
	v:=arg_v.Encode()
	if "GET"== arg_method||"get"==arg_method{
		arg_url+="?"+v
	}
	r, _ := http.NewRequest(arg_method, arg_url, ioutil.NopCloser(strings.NewReader(v)))
	if "POST" == arg_method {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}
	for _, v := range cookie {
		r.AddCookie(v)
	}



	for _, v := range cookie {
		r.AddCookie(v)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	re := ReData{}
	resp_body := w.Body.String()
	err := util.Json2S(resp_body, &re)
	if nil!=err {
		fmt.Println("json to struct failed--:%v", err)
		re.M = resp_body
	}
	return re
}


//--------------------------------old



//set session for test
func TestSetSession(hs *routing.HTTPSession) routing.HResult {
	uid, err := hs.R.Cookie("uid")
	uidInt64:=int64(0)
	if nil==err {
		uidInt64, _ = strconv.ParseInt(uid.Value, 10, 0)
	}
	username, err := hs.R.Cookie("username")
	if nil==err {
		hs.SetVal("username", username.Value)
	}
	hs.SetVal("uid", uidInt64)
	log.D("uid:%v,username:%v", hs.IntVal("uid"), hs.StrVal("username"))
	return routing.HRES_CONTINUE
}
func GetCartData(db *sql.DB, arg_sql string, arg_userId int64) ([]*RE_CART, error) {
	reCarts := []*RE_CART{}
	//这里可以加arg_sql的判断
	rows, err := db.Query(arg_sql, arg_userId)
	if err != nil {
		log.E(M1, err)
		return reCarts, errors.New(M1)
	}
	for rows.Next() {
		r := RE_CART{}
		err := rows.Scan(&r.ID, &r.USER_ID, &r.GOODS_ID, &r.GOODS_PRICE, &r.GOODS_NUM, &r.GOODS_IMG, &r.GOODS_NAME, &r.GOODS_URL,
			&r.SHOP_ID, &r.SHOP_NAME, &r.SHOP_URL, &r.SELLER_ID, &r.SELLER_NAME, &r.FROM_URL, &r.STATUS, &r.TIME)
		if err != nil {
			log.E(M4, err)
			return reCarts, errors.New(M4)
		}
		reCarts = append(reCarts, &r)
	}
	defer rows.Close();
	return reCarts, nil
}
func RequestUrl(db *sql.DB, arg_orderData ORDERS_ITEM) (string, string) {
	e := ""
	reUrl := ""
	callbackData := ORDERS_CALLBACK_DATA{arg_orderData.UsrId, arg_orderData.Status, arg_orderData.GoodsIds}
	v := url.Values{}
	//struct 到json str
	d, err := json.Marshal(callbackData);
	if err != nil {
		log.E("转Json失败！", err)
		return reUrl, e
	}
	fmt.Println(string(d))
	v.Set("data", string(d))
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	client := &http.Client{}
	Url := arg_orderData.FromUrl
	if ""==Url {
		log.E("请求路径为空！")
		e="请求路径为空！"
		return reUrl, e
	}
	req, _ := http.NewRequest("POST", Url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") //这个一定要加，不加form的值post不过去
	//	//添加token
	//	token, _ := hs.R.Cookie("token")
	//	req.AddCookie(token)
	resp, err := client.Do(req) //发送
	if nil!=err {
		log.E("请求失败，请重试！", err)
		e = "请求失败，请重试！"
		return reUrl, e
	}
	defer resp.Body.Close()     //一定要关闭resp.Body
	cbData, err := ioutil.ReadAll(resp.Body)
	if nil!=err {
		e = "请求错误！"
		log.E(e, err)
		return reUrl, e
	}
	getData := ReData{}
	if err := json.Unmarshal([]byte(cbData), &getData); err != nil {
		e = "传入数据有误！"
		log.E(e, err)
		return reUrl, e
	}
	fmt.Println(string(cbData), err)
	if getData.C!=0 {
		e = getData.M
		return reUrl, e
	}
	e = getData.M
	reUrl = getData.D.(string)
	fmt.Println(getData, e)
	return reUrl, e
}
func CallBack(db *sql.DB, arg_orderData []*ORDERS_ITEM) string {
	e := ""
	for _, orderData := range arg_orderData {
		//修改数据库同步状态
		//同步状态为1已同步
		syncData := 1
		_, e := RequestUrl(db, *orderData)
		if e!="" {
			log.E(e)
			return e
		}
		_, err := db.Exec(updateSync, syncData, orderData.OrderId)
		if nil != err {
			e = M1
			log.E(e, err)
			return e
		}
	}
	return e
}
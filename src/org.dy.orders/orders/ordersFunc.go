package orders
import (
	"database/sql"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"time"
	"net/url"
	"io/ioutil"
	"strings"
	"net/http"
	"strconv"
	"errors"
)

//check table
func CheckT(db *sql.DB) error {
	if nil==db {
		return errors.New("db conn fail")
	}
	if _, err := dbutil.DbQuery(db, chk_o_usr, 1); err != nil {
		log.D("o_usr table not found, auto creating...")
		err = dbutil.DbExecScript(db, o_usr)
		if err != nil {
			log.I("creat table o_usr failed--:%v", err)
			return err
		}
	}
	return nil
}
//add to cart example
func AddToCartTestF(hs *routing.HTTPSession, arg_idx  int64) ReData {
	log.D("AddToCartTestF r data--:%v", arg_idx)
	re := ReData{1, "", nil}

	mockCartData := []*O_CART{}
	cartData := O_CART{}
	cartData.USER_ID = 36
	cartData.GOODS_ID = 5
	cartData.GOODS_PRICE = 5
	cartData.GOODS_NUM = 1
	cartData.GOODS_IMG = "http://localhost:8082/resImg/1106100185/5_1438675199139484700_imgPre0.jpg"
	cartData.GOODS_NAME = "秋裙"
	cartData.GOODS_URL = "http://localhost:8082/resCenter.html#/Res/5"
	cartData.SHOP_ID = 1106100185
	cartData.SHOP_NAME =  "女娃娃"
	cartData.SHOP_URL = "http://localhost:8082/decorate/storeFirstPage.html?1106100185"
	cartData.SELLER_ID = 1111061127
	cartData.SELLER_NAME = "linj"
	cartData.FROM_URL = "http://localhost:8082/api/getUrl"
	cartData.STATUS = 20
	cartData_ := []*O_CART{}
	mockCartData = append(mockCartData, &cartData)

	cartData_ = append(cartData_, mockCartData[arg_idx])

	e, data := PostAddToCart(hs.R, cartData_)
	if e!="" {
		re.M = e
		return re
	}
	if len(data)>0 {
		re.D = data
	}
	re.C = 0
	return re
}
//post data to orders sys add to cart
func PostAddToCart(r *http.Request, arg_d []*O_CART) (string, []int64) {
	e := ""
	data := []int64{}
	log.D("PostAddToCart data--:%v", arg_d)
	e = checkPostAddToCart(arg_d)
	if e!="" {
		log.E(e)
		return e, data
	}
	//struct to json str
	d, err := json.Marshal(arg_d);
	if err != nil {
		e = "转Json失败！"
		log.E(e+"--:%v", err)
		return e, data
	}
	v := url.Values{}
	v.Set("data", string(d))
	//Encode form data
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	url := AddToOrdersUrl()
	if ""==url {
		log.E(M10)
		return e, data
	}
	req, _ := http.NewRequest("POST", url, body)
	//form post
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//add token
	token, _ := r.Cookie("token")
	req.AddCookie(token)
	//send http
	resp, err := client.Do(req)
	if nil!=err {
		e = "请求失败，请重试！"
		log.E(e+"--：%v", err)
		return e, data
	}
	//close resp.Body
	defer resp.Body.Close()

	cbData, err := ioutil.ReadAll(resp.Body)
	if nil!=err {
		e = "请求错误！"
		log.E(e+"--：%v", err)
		return e, data
	}
	getData := CbData{}
	if err := json.Unmarshal([]byte(cbData), &getData); err != nil {
		e = "返回数据有误！"
		log.E(e+"--：%v", err)
		return e, data
	}
	log.D("orders return data--:%v",getData)
	e = getData.Msg
	data = getData.Data
	return e, data
}
//check PostAddToCart data
func checkPostAddToCart(arg_d []*O_CART) string {
	e := ""
	for _, d := range arg_d {
		if d.USER_ID<1 {
			e=M6
			break
		}
		if d.GOODS_ID<1 {
			e=M6
			break
		}
		if d.GOODS_PRICE<0 {
			e=M6
			break
		}
		if d.GOODS_NUM<1 {
			e=M6
			break
		}
		if d.GOODS_IMG=="" {
			//			e=M6
			//			break
		}
		if d.GOODS_NAME=="" {
			e=M6
			break
		}
		if d.GOODS_URL=="" {
			e=M6
			break
		}
		if d.SHOP_ID<1 {
			//			e=M6
			//			break
		}
		if d.SHOP_NAME=="" {
			//			e=M6
			//			break
		}
		if d.SHOP_URL=="" {
			//			e=M6
			//			break
		}
		if d.SELLER_ID<1 {
			e=M6
			break
		}
		if d.SELLER_NAME=="" {
			e=M6
			break
		}
		if d.FROM_URL=="" {
			e=M6
			break
		}
		if d.STATUS<0 {
			//			e=M6
			//			break
		}
	}

	return e
}
func JsonRes(hs *routing.HTTPSession, code int64, data interface{}, msg string) routing.HResult {
	res := make(util.Map)
	res["code"] = code
	res["msg"] = msg
	res["data"] = data
	dBys, _ := json.Marshal(res)
	hs.W.Write(dBys)
	return routing.HRES_RETURN
}
//------------------OLD




func SyncOrdersF(db *sql.DB, arg_d G_ORDERS_DATA) string {
	log.D("arg_d:", arg_d)
	e := ""
	tx, _ := db.Begin()
	for i, _ := range arg_d.GoodsIds {
		if arg_d.GoodsIds[i]<1||arg_d.UsrId<1 {
			e = "传入数据有误"
			tx.Rollback()
			return e
		}
		//更新人员订单表
		stmt, err := tx.Exec(update_o_usr,
			arg_d.Status,
			arg_d.UsrId,
			arg_d.GoodsIds[i],
			arg_d.Token,
		)
		if nil!= err {
			log.E(M3, err)
			e = M3
			tx.Rollback()
			return e
		}
		temLen, err := stmt.RowsAffected()
		if nil!= err {
			log.E(M9, err)
			e = M9
			tx.Rollback()
			return e
		}
		if temLen<1 {
			log.E(M8, err)
			e = M8
			tx.Rollback()
			return e
		}
	}
	tx.Commit()
	return e
}
//获取请求路径
func GetUrlF(db *sql.DB, arg_d G_ORDERS_DATA) (string, string) {
	log.D("arg_d:", arg_d)
	e := ""
	token := ""
	token = uuid.New()
	tx, _ := db.Begin()
	for i, _ := range arg_d.GoodsIds {
		data := O_USR{}
		data.GOODS_ID=arg_d.GoodsIds[i]
		data.USR_ID = arg_d.UsrId
		data.TOKEN = token
		data.TIME = time.Now().Format("2006-01-02 15:04:05")
		if data.GOODS_ID<1||data.USR_ID<1 {
			e = "传入数据有误"
			tx.Rollback()
			return "", e
		}
		//插入人员订单表
		_, err := tx.Exec(ins_o_usr,
			data.USR_ID,
			data.GOODS_ID,
			data.STATUS,
			data.TOKEN,
			data.TIME,
		)
		if nil!= err {
			log.E(M3, err)
			e = M3
			tx.Rollback()
			return "", e
		}
	}
	tx.Commit()
	log.E("token:", token)
	return token, e
}



func AddToCartF(r *http.Request, db *sql.DB, arg_d []*G_ADD_CART, arg_uid int64) ReData {
	log.D("data:", arg_d)
	re := ReData{1, "", nil}
	//根据arg_d获取数据
	cartData := []*O_CART{}
	host := "http://"+r.Host
	getUrl := "/api/getUrl"
	for _, d := range arg_d {
		log.D("d:", d)
		selectData := ""
		selectData = selectGoodsInf
		rows, err := db.Query(selectData, d.GoodsId)
		if err != nil {
			log.E(M1, err)
			re.M = M1
			return re
		}
		for rows.Next() {
			g := O_CART_SQL{}
			err := rows.Scan(&g.GOODS_ID, &g.GOODS_PRICE, &g.GOODS_IMG, &g.GOODS_NAME,
				&g.SHOP_ID, &g.SHOP_NAME, &g.SELLER_ID, &g.SELLER_NAME)
			if err != nil {
				log.E(M4, err)
				re.M = M4
				return re
			}
			log.D("goods:", g)
			temPic := strings.Split(g.GOODS_IMG.String, "|")
			// temOI[0]订单数据 temOI[1]地址ID temOI[2]商品ID 数量
			log.D("len:", len(temPic), temPic[0])


			cartData_ := O_CART{}
			cartData_.USER_ID = arg_uid
			cartData_.GOODS_ID = d.GoodsId
			cartData_.GOODS_PRICE = g.GOODS_PRICE.Float64
			cartData_.GOODS_NUM = d.GoodsNum
			cartData_.GOODS_IMG = host+"/"+temPic[0]
			cartData_.GOODS_NAME = g.GOODS_NAME.String
			cartData_.GOODS_URL = host+"/resCenter.html#/Res/"+strconv.FormatInt(g.GOODS_ID.Int64, 10)
			cartData_.SHOP_ID = g.SHOP_ID.Int64
			cartData_.SHOP_NAME =  g.SHOP_NAME.String
			cartData_.SHOP_URL = host+"/decorate/storeFirstPage.html?"+strconv.FormatInt(g.SHOP_ID.Int64, 10)
			cartData_.SELLER_ID = g.SELLER_ID.Int64
			cartData_.SELLER_NAME = g.SELLER_NAME.String
			cartData_.FROM_URL = host+getUrl
			cartData_.STATUS = d.Status
			log.D("---------------------------\n\n\n\ndata:", cartData_)
			cartData = append(cartData, &cartData_)
		}
		defer rows.Close()
	}
	e, data := PostAddToCart(r, cartData)
	if e!="" {
		re.M = e
		return re
	}
	if len(data)>0 {
		re.D = data
	}
	re.C = 0
	return re
}

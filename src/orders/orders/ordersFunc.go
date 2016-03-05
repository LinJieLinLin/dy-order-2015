package orders
import (
	C "common"
	D "databaseConn"
	"github.com/Centny/gwf/log"
	"strconv"
	"math/rand"
	"time"
	"orders/cart"
	"database/sql"
)




//---------------old
//提交订单
func SubmitOrderF(arg_d []*C.G_SubmitData, arg_userId int64) C.ReData {
	log.D("data:", arg_d)
	re := C.ReData{1, "", nil}
	re.M = checkSubmitOrderF(arg_d, arg_userId)
	if re.M!="" {
		return re
	}

	db, err := D.GetConn()
	if nil!=err {
		log.E(C.M1, err)
		re.M = C.M1
		return re
	}
	tx, _ := db.Begin()
	ordersCreateTime := time.Now().Format("2006-01-02 15:04:05")
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderNo := time.Now().Format("200601")+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+
	strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+
	strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))
	pubOrderNo := "p"+orderNo
	reData := C.RE_ORDERS{pubOrderNo}

	allPrice := float64(0)
	for _, d := range arg_d {
		cartData := []*C.RE_CART{}
		orderPrice := float64(0)
		//orderType订单类型(0WEB订单，10手机订单)
		orderType := int64(0)
		//		if reData.IsOne==0 {
		//			reData.IsOne = 1
		//		}
		for _, cartId := range d.CartId {
			log.D("cartId", cartId)
			reCarts := []*C.RE_CART{}
			getCartSql := getCart
			getCartSql+= " and (STATUS = 0 or STATUS= 20 ) and id="+strconv.FormatInt(cartId, 10)
			log.D("sql:", getCartSql)
			reCarts, err=C.GetCartData(db, getCartSql, arg_userId)
			if err!=nil {
				log.E(C.M1, err)
				re.M = C.M1
				tx.Rollback()
				return re
			}
			if len(reCarts)==0 {
				log.E(C.M7, err)
				re.M = C.M7
				tx.Rollback()
				return re
			}
			//读取购物车数据
			orderPrice = orderPrice+reCarts[0].GOODS_PRICE*float64(reCarts[0].GOODS_NUM)
			cartData = append(cartData, reCarts[0])
		}
		//这里写订单号重复的判断
		//这里写订单号重复的判断
		//这里写订单号重复的判断

		//订单号：年月+9位随机数
		ra := rand.New(rand.NewSource(time.Now().UnixNano()))
		orderNo := time.Now().Format("200601")+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+
		strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+
		strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))+strconv.Itoa(ra.Intn(10))
		log.D("订单号：", orderNo)
		ordersD := C.O_ORDERS{
			0,
			orderNo,
			pubOrderNo,
			cartData[0].USER_ID,
			cartData[0].SHOP_ID,
			cartData[0].SHOP_NAME,
			cartData[0].SHOP_URL,
			cartData[0].SELLER_ID,
			cartData[0].SELLER_NAME,
			orderType,
			0,
			cartData[0].FROM_URL,
			orderPrice,
			ordersCreateTime,
			"",
			ordersCreateTime,
			0,
		}
		allPrice+=ordersD.PRICE
		//请求子系统生成token
		ordersData := []*C.ORDERS_ITEM{}
		ordersOneData := C.ORDERS_ITEM{}
		ordersOneData.UsrId = cartData[0].USER_ID
		ordersOneData.FromUrl = cartData[0].FROM_URL
		for i, _ := range cartData {
			ordersOneData.GoodsIds = append(ordersOneData.GoodsIds, cartData[i].GOODS_ID)
		}
		ordersData = append(ordersData, &ordersOneData)
		url, e := checkUrl(db, ordersData)
		if e!="" {
			log.E(C.M3, e)
			re.M = e
			tx.Rollback()
			return re
		}
		ordersD.FROM_URL = url
		//判断是否有免费的商品
		if ordersD.PRICE==0 {
			ordersD.STATUS=70
			ordersD.SYNC = 1
			//			reData.OrderNo = "free"
		}
		log.D("订单数据：", ordersD)
		//插入订单表
		stmt, err := tx.Exec(insertOrders,
			ordersD.ORD_NO,
			ordersD.PUB_ORD_NO,
			ordersD.USR_ID,
			ordersD.SHOP_ID,
			ordersD.SHOP_NAME,
			ordersD.SHOP_URL,
			ordersD.SELLER_ID,
			ordersD.SELLER_NAME,
			ordersD.ORD_TYPE,
			ordersD.STATUS,
			ordersD.FROM_URL,
			ordersD.PRICE,
			ordersD.CREATE_TIME,
			ordersD.TIME,
			ordersD.SYNC,
		)
		if nil!= err {
			log.E(C.M3, err)
			re.M = C.M3
			tx.Rollback()
			return re
		}
		ordersId, err := stmt.LastInsertId()
		if nil != err {
			log.E(C.M3, err)
			re.M = C.M3
			tx.Rollback()
			return re
		}

		ordersSnyc := C.ORDERS_ITEM{}
		for _, orderItem := range cartData {
			if ordersD.SELLER_ID!=orderItem.SELLER_ID {
				log.E(C.M6, err)
				re.M = C.M6
				tx.Rollback()
				return re
			}
			log.D("订单ID", ordersId, orderItem)
			//插入订单详情
			_, err := tx.Exec(insertOrdersItem,
				ordersId,
				orderItem.GOODS_PRICE,
				orderItem.GOODS_ID,
				orderItem.GOODS_NUM,
				orderItem.GOODS_IMG,
				orderItem.GOODS_NAME,
				orderItem.GOODS_URL,
				0,
				orderItem.TIME,
			)
			if nil!= err {
				log.E(C.M3, err)
				re.M = C.M3
				tx.Rollback()
				return re
			}
			//更改购物车
			data := []*C.G_CartData{}
			data_ := C.G_CartData{
				orderItem.ID,
				1,
				0,
			}
			if orderItem.STATUS==0 {
				//下单完成删除购物车数据
				data_.Status=-1
			}else if orderItem.STATUS==20 {
				data_.Status=-1
			}
			if ordersD.STATUS==70 {
				ordersSnyc.GoodsIds = append(ordersSnyc.GoodsIds, orderItem.GOODS_ID)
			}
			data = append(data, &data_)
			re = cart.EditCartF(data, arg_userId, tx)
		}
		if ordersD.STATUS==70 {
			ordersSnyc.OrderId = ordersId
			ordersSnyc.UsrId = ordersD.USR_ID
			ordersSnyc.Status = 10
			ordersSnyc.FromUrl = ordersD.FROM_URL
			_, e := C.RequestUrl(db, ordersSnyc)
			if e!="" {
				log.E(C.M3, e)
				re.M = C.M3
				tx.Rollback()
				return re
			}
		}
	}
	if allPrice==0 {
		reData.OrderNo = "free"
	}
	tx.Commit()
	re.C = 0
	re.D = reData
	return re
}
func checkSubmitOrderF(arg_d []*C.G_SubmitData, arg_userId int64) string {
	err := ""
	if len(arg_d)==0 {
		log.E(C.M5)
		err = C.M5
	}
	for _, d := range arg_d {
		if len(d.CartId)==0 {
			log.E(C.M5)
			err = C.M5
			break
		}
		for _, temD := range d.CartId {
			if temD<1 {
				err = "传入数据有误！"
				log.E(err)
				break
			}
		}
	}
	if arg_userId<1 {
		err = C.M6
	}
	return err
}
func checkUrl(db *sql.DB, arg_orderData []*C.ORDERS_ITEM) (string, string) {
	e := ""
	url := ""
	for _, orderData := range arg_orderData {
		url, e=C.RequestUrl(db, *orderData)
		if e!="" {
			return url, e
		}
	}
	log.I("url:", url)
	return url, e
}
//获取订单数据
func GetOrdersDataF(arg_d C.G_OrdersData) C.ReData {
	log.D("data:", arg_d)
	re := C.ReData{1, "", nil}
	re.M = checkGetOrdersDataF(arg_d)
	if re.M!="" {
		return re
	}

	db, err := D.GetConn()
	if nil!=err {
		log.E(C.M1, err)
		re.M = C.M1
		return re
	}
	tx, _ := db.Begin()

	reData := C.RE_ORDERS_DATA{}
	//总数量/每页数量 为页数
	dataSum := int64(0)
	//总页数
	pageSum := int64(1)
	//获取总页数
	getOrdersNum := getDataSum_fromUrl
	if arg_d.Status!=-1 {
		getOrdersNum = getDataSum_fromUrlByStatus
	}
	log.D("getOrdersNum:", getOrdersNum, arg_d.UserId, arg_d.FromUrl, arg_d.Status)
	if err := db.QueryRow(getOrdersNum, arg_d.UserId, "%"+arg_d.FromUrl+"%", arg_d.Status).Scan(&dataSum); err != nil {
		log.E(C.M1, err)
		re.M = C.M1
		tx.Rollback()
		return re
	}
	if dataSum==0 {
		re.C = 0
		reData.Page = 1
		re.D = reData
		re.M = C.M10
		tx.Rollback()
		return re
	}
	pageSum=getPageNum(dataSum, arg_d.PSize)
	if arg_d.PNum>pageSum {
		log.E(C.M9)
		re.M = C.M9
		tx.Rollback()
		return re
	}
	reData.Page = pageSum
	getOrdersList := getOrdersListByFromUrl
	if arg_d.Status!=-1 {
		getOrdersList=getOrdersListByStatus
	}
	log.D("getOrdersList:", getOrdersList, arg_d.UserId, "%"+arg_d.FromUrl+"%", arg_d.Status, (arg_d.PNum-1)*arg_d.PSize, arg_d.PSize)
	rows, err := db.Query(getOrdersList, arg_d.UserId, "%"+arg_d.FromUrl+"%", arg_d.Status, (arg_d.PNum-1)*arg_d.PSize, arg_d.PSize)
	if err != nil {
		log.E(C.M1, err)
		re.M = C.M1
		tx.Rollback()
		return re
	}
	for rows.Next() {
		s := C.ORDERS{}
		err := rows.Scan(&s.ORD_ID, &s.ORD_NO, &s.SHOP_ID, &s.SHOP_NAME, &s.SHOP_URL, &s.SELLER_ID, &s.SELLER_NAME, &s.ORD_TYPE, &s.STATUS, &s.FROM_URL, &s.PRICE, &s.CREATE_TIME)
		if err != nil {
			log.E(C.M4, err)
			re.M = C.M4
			tx.Rollback()
			return re
		}
		goods := []*C.GOODS{}
		getGoodsUrl := goodsUrl
		log.D("getGoodsUrl:", getGoodsUrl)
		rows, err := db.Query(getGoodsUrl, s.ORD_ID)
		if err != nil {
			log.E(C.M1, err)
			re.M = C.M1
			tx.Rollback()
			return re
		}
		for rows.Next() {
			g := C.GOODS{}
			err := rows.Scan(&g.ID, &g.GOODS_ID, &g.GOODS_PRICE, &g.GOODS_NUM, &g.GOODS_IMG, &g.GOODS_NAME, &g.GOODS_URL)
			if err != nil {
				log.E(C.M4, err)
				re.M = C.M4
				tx.Rollback()
				return re
			}
			log.D("goods:", g)
			goods = append(goods, &g)
		}
		defer rows.Close();
		log.D("ORDERS:", s)
		s.GOODS = goods
		reData.OrdersList = append(reData.OrdersList, &s)
	}
	defer rows.Close();
	tx.Commit()
	re.C = 0
	re.D = reData
	return re
}
func checkGetOrdersDataF(arg_d C.G_OrdersData) string {
	err := ""
	if arg_d.UserId<1 {
		err = "用户ID有误"
	}else if arg_d.PNum<1 {
		err = "页数有误"
	}else if arg_d.PSize<1 {
		err = "每页数量有误"
	}
	log.D(err)
	return err
}
func getPageNum(arg_sum, arg_pSize int64) int64 {
	page := int64(0)
	if arg_sum%arg_pSize == 0 {
		page = arg_sum/arg_pSize
	}else if arg_sum%arg_pSize  != 0 {
		page = arg_sum/arg_pSize+1
	}
	return page
}
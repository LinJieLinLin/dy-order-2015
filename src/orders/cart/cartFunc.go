package cart
import (
	C "common"
	D "databaseConn"
	"github.com/Centny/gwf/log"
	"time"
	"database/sql"
	"strconv"
	"errors"
	"github.com/Centny/gwf/dbutil"
)

//add to cart by arg_d list
func AddToCartF(arg_d []*C.O_CART) C.ReData {
	re := C.ReData{1, "", nil}
	log.D("AddToCartF input data len--:%v", len(arg_d))
	err := checkAddToCartF(arg_d)
	if nil!=C.ReFuncErr(&re, err, C.M6) {
		return re
	}

	db, err := D.GetConn()
	if nil!=C.ReFuncErr(&re, err, C.M1) {
		return re
	}
	tx, _ := db.Begin()
	cartIds := []int64{}
	for _, d := range arg_d {
		id := int64(0)
		re, id=AddToCartByOne(tx, db, d)
		if re.C!=0 {
			tx.Rollback()
			return re
		}
		if id>0{
			cartIds = append(cartIds, id)
		}
	}
	tx.Commit()
	re.C = 0
	re.D = cartIds
	return re
}
// add to cart by one
func AddToCartByOne(tx *sql.Tx, db *sql.DB, arg_d *C.O_CART) (C.ReData, int64) {
	log.D("AddToCartByOne input data--:%v", arg_d)
	re := C.ReData{1, "", nil}
	id := int64(0)

	err := checkAddToCartFByOne(arg_d)
	if nil!=C.ReFuncErr(&re, err, C.M6) {
		return re, id
	}

	switch arg_d.STATUS {
	case 0:
		//check goods in cart
		temCount := 0
		log.D("check goods in o_cart,sql--:%v--:%v--:%v--:%v", selectCart, arg_d.USER_ID, arg_d.FROM_URL, arg_d.GOODS_ID)
		err := db.QueryRow(selectCart, arg_d.USER_ID, arg_d.FROM_URL, arg_d.GOODS_ID).Scan(&temCount);
		if nil!=C.ReFuncErrTx(tx, &re, err, C.M1) {
			return re, id
		}
		if temCount!=0 {
			//update cart
			log.D("update goods num sql--:%v--:%v--:%v--:%v--:%v--:%v--:%v", updateCart, arg_d.GOODS_PRICE, arg_d.GOODS_NUM, arg_d.GOODS_IMG, arg_d.USER_ID, arg_d.FROM_URL, arg_d.GOODS_ID)
			count, err := dbutil.DbUpdate2(tx, updateCart, arg_d.GOODS_PRICE, arg_d.GOODS_NUM, arg_d.GOODS_IMG, arg_d.USER_ID, arg_d.FROM_URL, arg_d.GOODS_ID)
			if nil!=C.ReFuncErrTx(tx, &re, err, C.M1) {
				return re, id
			}
			log.I("update num is--:%v", count)
			re.C = 0
			return re,id
		}
	case 20:
		//delete Direct data and insert new goods data
		log.D("del Cart data Sql--:%v--:%v--:%v--:%v", delCartSql, arg_d.STATUS, arg_d.USER_ID, "%"+arg_d.FROM_URL+"%")
		count, err := dbutil.DbUpdate2(tx, delCartSql, arg_d.STATUS, arg_d.USER_ID, "%"+arg_d.FROM_URL+"%")
		if nil!=C.ReFuncErrTx(tx, &re, err, C.M1) {
			return re, id
		}
		log.I("del num is--:%v", count)
	}
	//insert
	arg_d.TIME = time.Now().Format("2006-01-02 15:04:05")
	log.D("insert Cart data sql--:%v", insertCart, arg_d.USER_ID, arg_d.GOODS_ID, arg_d.GOODS_PRICE, arg_d.GOODS_NUM, arg_d.GOODS_IMG,
		arg_d.GOODS_NAME, arg_d.GOODS_URL, arg_d.SHOP_ID, arg_d.SHOP_NAME, arg_d.SHOP_URL, arg_d.SELLER_ID, arg_d.SELLER_NAME,
		arg_d.FROM_URL, arg_d.STATUS, arg_d.TIME, arg_d.ADD1,
	)
	temId, err := dbutil.DbInsert2(tx, insertCart, arg_d.USER_ID, arg_d.GOODS_ID, arg_d.GOODS_PRICE, arg_d.GOODS_NUM, arg_d.GOODS_IMG, arg_d.GOODS_NAME,
		arg_d.GOODS_URL, arg_d.SHOP_ID, arg_d.SHOP_NAME, arg_d.SHOP_URL, arg_d.SELLER_ID, arg_d.SELLER_NAME, arg_d.FROM_URL, arg_d.STATUS, arg_d.TIME, arg_d.ADD1,
	)
	if nil!=C.ReFuncErrTx(tx, &re, err, C.M3) {
		return re, id
	}
	if temId==0{
		re.M = C.M3
		return re, id
	}
	id = temId
	re.C = 0
	return re, id
}
//check AddToCartF data
func checkAddToCartF(arg_Data []*C.O_CART) error {
	log.D("checkAddToCartF input data len--:%v", len(arg_Data))
	e := ""
	if len(arg_Data)==0 {
		e=C.M5
		log.E("checkAddToCartF --:%v", e)
		return errors.New(e)
	}
	for _, data := range arg_Data {
		err := checkAddToCartFByOne(data)
		if nil!=err {
			return err
		}
	}
	return nil
}
func checkAddToCartFByOne(arg_Data *C.O_CART) error {
	//	return err
	log.D("checkAddToCartFByOne data--:%v", arg_Data)
	e := ""
	if arg_Data.USER_ID<1 {
		e = "用户ID有误！"
	}
	if arg_Data.GOODS_ID<1 {
		e = "商品ID有误！"
	}
	if arg_Data.GOODS_PRICE<0 {
		e = "商品价格有误！"
	}
	if arg_Data.GOODS_NUM<1 {
		e = "商品数量有误！"
	}
	//	if arg_Data.GOODS_IMG=="" {
	//		err = "商品图片路径有误！"
	//		break
	//	}
	if arg_Data.GOODS_NAME=="" {
		e = "商品名称有误！"
	}
	if arg_Data.GOODS_URL=="" {
		e = "商品链接有误！"
	}
	if arg_Data.SHOP_ID<0 {
		e = "商店ID有误！"
	}
	if arg_Data.SELLER_ID<1 {
		e = "卖家ID有误！"
	}
	if arg_Data.SELLER_NAME=="" {
		e = "卖家有误！"
	}
	if arg_Data.FROM_URL=="" {
		e = "来源URL有误！"
	}
	if arg_Data.STATUS!=0&&arg_Data.STATUS!=20 {
		e = "商品状态有误！"
	}
	if ""!=e {
		return errors.New(e)
	}
	return nil
}

//get cart data
func GetCartF(arg_userId int64, arg_d C.G_CartList) C.ReData {
	re := C.ReData{1, "", nil}
	log.D("GetCartF input data arg_userId--:%v arg_d--:%v", arg_userId, arg_d)
	err := checkGetCartF(arg_userId, arg_d)
	if nil!=C.ReFuncErr(&re, err, C.M6) {
		return re
	}

	db, err := D.GetConn()
	if nil!=C.ReFuncErr(&re, err, C.M1) {
		return re
	}
	status := int64(0)
	temL := len(arg_d.CartId)
	if temL>0 {
		status = 20
	}
	temCartIds := "0"
	for i := 0; i<temL; i++ {
		temCartIds+=","+strconv.FormatInt(arg_d.CartId[i], 10)
	}
	reCartData := []*C.SELLERS{}
	//get shop info
	selectShop := selectShopInf
	switch status {
	case 0:
		selectShop+=" ORDER BY ID DESC"
	case 20:
		selectShop+=" and ID in("+temCartIds+") ORDER BY ID DESC"
	}
	log.D("select Shop data from o_cart sql --:%v--:%v--:%v--:%v", selectShop, arg_userId, "%"+arg_d.FromUrl+"%", temCartIds)
	err = dbutil.DbQueryS(db, &reCartData, selectShop, arg_userId, "%"+arg_d.FromUrl+"%", status)
	if nil!=C.ReFuncErr(&re, err, C.M1) {
		return re
	}
	log.I("shop Data Len:--:%v", len(reCartData))
	if len(reCartData)<1 {
		re.C = 0
		re.M = "购物车为空！"
		return re
	}
	//get goods info
	reGoods := []*C.GOODS{}
	selectGoods := selectGoodsInf
	if status==20 {
		selectGoods+="and ID in("+temCartIds+") ORDER BY ID DESC"
	}
	log.D("select goods data from o_cart sql--:%v--:%v--:%v", selectGoods, arg_userId, status)
	err = dbutil.DbQueryS(db, &reGoods, selectGoods, arg_userId, status)
	if nil!=C.ReFuncErr(&re, err, C.M1) {
		return re
	}
	log.I("Goods Data Len:--:%v", len(reGoods))
	//add goodsData to cartData by SELLER_ID
	for _, goods := range reGoods {
		for _, cartData := range reCartData {
			if cartData.SELLER_ID==goods.SELLER_ID {
				log.D("select goods data is--:%v", goods)
				cartData.GOODS = append(cartData.GOODS, goods)
			}
		}
	}

	re.C = 0
	re.M = "购物车为空！"
	if len(reCartData)>0 {
		re.D = reCartData
		re.M = ""
	}
	return re
}
//check GetCartF data
func checkGetCartF(arg_userId int64, arg_d C.G_CartList) error {
	log.D("checkGetCartF input data --%:v--:%v",arg_userId,arg_d)
	err := errors.New(C.M6)
	if arg_userId<1 {
		log.E(C.M6+"arg_userId--:%v", arg_userId)
		return err
	}
	for _, d := range arg_d.CartId {
		if d<1 {
			log.E(C.M6+"CartId--:%v", d)
			return err
		}
	}
	return nil
}

//edit cart data by list
func EditCartF(arg_d []*C.G_CartData, arg_userId int64, arg_tx *sql.Tx) C.ReData {
	log.D("EditCartF input data len--:%v arg_userId--:%v",arg_d,arg_userId)
	re := C.ReData{1, "", nil}
	err := checkEditCartF(arg_d, arg_userId);
	if nil!=C.ReFuncErr(&re, err, C.M6) {
		return re
	}

	db, err := D.GetConn()
	if nil!=C.ReFuncErr(&re, err, C.M1) {
		return re
	}

	var tx *sql.Tx
	if nil==arg_tx {
		tx, _ = db.Begin()
	}else if nil!=arg_tx {
		tx=arg_tx
	}

	for _, d := range arg_d {
		re=EditCartFByOne(tx, d, arg_userId)
		if 1==re.C {
			return re
		}
	}

	if nil==arg_tx {
		tx.Commit()
	}
	re.C = 0
	return re
}
//edit cart data by one
func EditCartFByOne(arg_tx *sql.Tx, arg_d *C.G_CartData, arg_userId int64) C.ReData {
	log.D("EditCartFByOne input data--:%v arg_userId--:%v", arg_d, arg_userId)
	re := C.ReData{1, "", nil}
	err := checkEditCartFByOne(arg_d, arg_userId);
	if nil!=C.ReFuncErr(&re, err,C.M6) {
		return re
	}
	editSql := ""
	v := arg_d.Status
	switch arg_d.Status {
	//from status 30 to 0
	case -3:
		editSql = editCartByUpdateStatus
		v=0
	case -1:
		editSql = editCartByDelete
	case 0:
		editSql = editCartByUpdateGoodsNum
		v = arg_d.Num
	case 30:
		editSql = editCartByUpdateStatus
	}

	log.D("edit 0_cart Sql--:%v--:%v--:%v--:%v", editSql, v, arg_d.CartId, arg_userId)
	count, err := dbutil.DbUpdate2(arg_tx, editSql, v, arg_d.CartId, arg_userId)
	if nil!=C.ReFuncErrTx(arg_tx, &re, err, C.M1) {
		return re
	}
	log.I("edit num--:%v", count)
	//update or delete num
	if count<1 {
		log.E(C.M8, err)
		re.M = C.M8
		arg_tx.Rollback()
		return re
	}
	re.C = 0
	return re
}
//checkEditCartF data by list
func checkEditCartF(arg_d []*C.G_CartData, arg_userId int64) error {
	log.D("checkEditCartF input data len--：%v--：%v",len(arg_d),arg_userId)
	e := ""
	if len(arg_d)==0 {
		e=C.M5
		return errors.New(e)
	}
	for _, d := range arg_d {
		err := checkEditCartFByOne(d, arg_userId)
		if nil!=err {
			return err
		}
	}
	return nil
}
//checkEditCartF data by one
func checkEditCartFByOne(arg_d *C.G_CartData, arg_userId int64) error {
	log.D("checkEditCartF input data --：%v--：%v",arg_d,arg_userId)
	e := ""
	if arg_d.CartId<1 {
		e = "购物车ID有误！"
	}
	if arg_d.Num<1 {
		e = "商品数量有误！"
	}
	if arg_d.Status!=0&&arg_d.Status!=30&&arg_d.Status!=-1&&arg_d.Status!=-3 {
		e = "购物车状态有误！"
	}
	if arg_userId<1 {
		e = "用户ID有误！"
	}
	if ""!=e {
		return errors.New(e)
	}
	return nil
}




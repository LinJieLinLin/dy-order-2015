package alipay
import (
	C "common"
	D "databaseConn"
	"github.com/Centny/gwf/log"
	"database/sql"
)

func GetOrdersPrice(arg_orderNo string) (float64, string) {
	log.D("GetOrdersPrice_data:", arg_orderNo)
	allPrice := float64(0)
	e := ""
	if arg_orderNo=="" {
		e= "订单号为空！"
		return allPrice, e
	}
	db, err := D.GetConn()
	if nil!=err {
		log.E(C.M1, err)
		e = C.M1
		return allPrice, e
	}
	allStatus := int64(0)
	var allPrice_ sql.NullFloat64
	var allStatus_ sql.NullInt64
	var selectOrdersSql = selectOrdersPrice
	//判断是否为下单时的公共订单号带P开头为，否则为单一的订单
	if string(arg_orderNo[0])=="p"{
		selectOrdersSql = selectOrdersPriceP
	}
	if err := db.QueryRow(selectOrdersSql, arg_orderNo).Scan(&allPrice_, &allStatus_); err != nil {
		log.E(C.M1, err)
		e = C.M1
		return allPrice, e
	}
	allPrice = allPrice_.Float64
	allStatus = allStatus_.Int64
	if allStatus>0 {
		e="订单已支付！"
		log.E(e,allStatus)
	}
	if allPrice<=0 {
		e = "订单号不存在！"
	}
	return allPrice, e
}
//支付成功修改订单状态
func EditOrders(arg_orderNo string) string {
	e := ""
	_, e = GetOrdersPrice(arg_orderNo)
	if ""!=e {
		if e=="订单已支付！"{
			er:=CheckSync(arg_orderNo)
			if ""!=er{
				e = er
			}
		}
		return e
	}
	db, err := D.GetConn()
	if nil!=err {
		log.E(C.M1, err)
		e = C.M1
		return e
	}
	_,err=db.Exec(updateOrderStatus,arg_orderNo)
	if nil != err {
		log.E(C.M1, err)
		e = C.M1
		return e
	}
	//回调请求
	e=CheckSync(arg_orderNo)
	return e
}
func CheckSync(arg_orderNo string) string{
	e := ""
	db, err := D.GetConn()
	if nil!=err {
		log.E(C.M1, err)
		e = C.M1
		return e
	}
	var selectOrdersIdsSql = selectOrdersIds
	//判断是否为下单时的公共订单号带P开头为，否则为单一的订单
	if string(arg_orderNo[0])=="p"{
		selectOrdersIdsSql = selectOrdersIdsP
	}
	rows, err := db.Query(selectOrdersIdsSql, arg_orderNo)
	if err != nil {
		log.E(C.M1, err)
		e = C.M1
		return e
	}
	ordersId := []*C.ORDERS_ITEM{}
	for rows.Next() {
		temOrdersId := C.ORDERS_ITEM{}
		var ordersId_ sql.NullInt64
		var fromUrl_ sql.NullString
		var usrId_ sql.NullInt64

		err := rows.Scan(&ordersId_,&fromUrl_,&usrId_)
		if err != nil {
			log.E(C.M4, err)
			e = C.M4
			return e
		}
		goodsIds := []int64{}
		rows, err := db.Query(selectOrdersGoodsId, ordersId_.Int64)
		if err != nil {
			log.E(C.M1, err)
			e = C.M1
			return e
		}
		for rows.Next() {
			var goodsId_ sql.NullInt64
			err := rows.Scan(&goodsId_)
			if err != nil {
				log.E(C.M4, err)
				e = C.M4
				return e
			}
			log.D("goodId:",goodsId_)
			goodsIds  = append(goodsIds, goodsId_.Int64)
		}
		defer rows.Close();

		temOrdersId.OrderId = ordersId_.Int64
		temOrdersId.GoodsIds = goodsIds
		temOrdersId.UsrId = usrId_.Int64
		temOrdersId.Status = 10
		temOrdersId.FromUrl = fromUrl_.String
		ordersId = append(ordersId, &temOrdersId)
	}
	defer rows.Close();
	log.I("订单ID:",ordersId)
	e=C.CallBack(db,ordersId)
	if e!=""{
		log.E("同步结果:",e)
		e = "同步失败！"
	}
	return e
}

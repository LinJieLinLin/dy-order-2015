package alipay
const(
	selectOrdersPriceP=`
	SELECT SUM(b.GOODS_NUM*b.GOODS_PRICE),sum(a.STATUS)
	from o_orders a,o_orders_item b WHERE
	a.ID=b.ORD_ID and a.PUB_ORD_NO = ? and a.STATUS<>70
	`
	selectOrdersPrice=`
	SELECT SUM(b.GOODS_NUM*b.GOODS_PRICE),sum(a.STATUS)
	from o_orders a,o_orders_item b WHERE
	a.ID=b.ORD_ID and a.ORD_NO = ? and a.STATUS<>70
	`
	selectOrdersIdsP=`
	SELECT a.ID,a.FROM_URL,a.USR_ID
	from o_orders a WHERE
	a.PUB_ORD_NO =? and a.STATUS<>70 and SYNC<>1
	`
	selectOrdersIds=`
	SELECT a.ID,a.FROM_URL,a.USR_ID
	from o_orders a WHERE
	a.ORD_NO =? and a.STATUS<>70 and SYNC<>1
	`
	selectOrdersGoodsId = `
	SELECT a.GOODS_ID
	from o_orders_item a WHERE
	a.ORD_ID =? and a.STATUS=0
	`

	updateOrderStatus=`
	update o_orders set STATUS = 10 WHERE
	PUB_ORD_NO=? and STATUS<>70
	`
)
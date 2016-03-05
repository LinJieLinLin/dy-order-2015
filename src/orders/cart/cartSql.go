package cart

const (
	selectCart = `
	SELECT COUNT(ID) FROM o_cart WHERE USER_ID = ? and FROM_URL = ? and GOODS_ID = ? and STATUS = 0
	`
	insertCart = "INSERT INTO o_cart ("+
	"`USER_ID`, `GOODS_ID`, `GOODS_PRICE`, `GOODS_NUM`, `GOODS_IMG`, `GOODS_NAME`, `GOODS_URL`, `SHOP_ID`, "+
	"`SHOP_NAME`, `SHOP_URL`, `SELLER_ID`, `SELLER_NAME`, `FROM_URL`,`STATUS`, `TIME`, `ADD1`"+
	") VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	updateCart = `
	update o_cart set GOODS_PRICE=?,GOODS_NUM=GOODS_NUM+?,GOODS_IMG=? where USER_ID = ? and FROM_URL = ? and GOODS_ID = ? and STATUS = 0
	`

	getCart = `
	select
	ID,
	USER_ID,
	GOODS_ID,
	GOODS_PRICE,
	GOODS_NUM,
	GOODS_IMG,
	GOODS_NAME,
	GOODS_URL,
	SHOP_ID,
	SHOP_NAME,
	SHOP_URL,
	SELLER_ID,
	SELLER_NAME,
	FROM_URL,
	TIME
 	from o_cart WHERE USER_ID=?
	`

	selectShopInf = `select DISTINCT(SELLER_ID),SELLER_NAME,SHOP_ID,SHOP_NAME,SHOP_URL,FROM_URL from o_cart WHERE USER_ID=?
	and FROM_URL LIKE ? and STATUS=?`
	selectGoodsInf = `
	select ID,USER_ID,SELLER_ID,GOODS_ID,GOODS_PRICE,GOODS_NUM,GOODS_IMG,GOODS_NAME,GOODS_URL,TIME from o_cart WHERE USER_ID=? and STATUS=?
	`
	editCartByUpdateGoodsNum = `update o_cart set GOODS_NUM=? where 1=1 and id=? and user_id = ? and status<>20
	`
	editCartByUpdateStatus = `update o_cart set STATUS=? where 1=1 and id=? and user_id = ? and status<>20
	`
	editCartByDelete = `delete from o_cart where STATUS<>? and id=? and user_id = ?
	`
	delCartSql = `
	delete from o_cart where STATUS=? and user_id = ? and FROM_URL LIKE ?
	`
)
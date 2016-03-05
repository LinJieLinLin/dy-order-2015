package orders
import (
	"database/sql"
)
const (
	M1 = "数据库连接出错。"
	M2 = "读取用户信息失败。"
	M3 = "插入数据失败。"
	M4 = "读取数据失败。"
	M5 = "传入数据为空。"
	M6 = "传入数据有误。"
	M7 = "订单号不存在。"
	M8 = "无数据更新。"
	M9 = "更新数据失败。"
	M10 = "读取配置文件出错！"
//sql


)
//table Data
type O_USR struct {
	ID       int64
	USR_ID   int64
	GOODS_ID int64
	STATUS   int64
	TOKEN    string
	TIME     string
	//	ADD1	string
	//	ADD2	string
}
type O_CART struct {
	//	ID          int64
	USER_ID     int64
	GOODS_ID    int64
	GOODS_PRICE float64
	GOODS_NUM   int64
	GOODS_IMG   string
	GOODS_NAME  string
	GOODS_URL   string
	SHOP_ID     int64
	SHOP_NAME   string
	SHOP_URL    string
	SELLER_ID   int64
	SELLER_NAME string
	FROM_URL    string
	STATUS      int64
	//	TIME        string
	//	ADD1        string
	//	ADD2        string
}
type O_CART_SQL struct {
	//	ID      sql.NullInt64
	USER_ID     sql.NullInt64
	GOODS_ID    sql.NullInt64
	GOODS_PRICE sql.NullFloat64
	GOODS_NUM   sql.NullInt64
	GOODS_IMG   sql.NullString
	GOODS_NAME  sql.NullString
	GOODS_URL   sql.NullString
	SHOP_ID     sql.NullInt64
	SHOP_NAME   sql.NullString
	SHOP_URL    sql.NullString
	SELLER_ID   sql.NullInt64
	SELLER_NAME sql.NullString
	FROM_URL    sql.NullString
	STATUS      sql.NullInt64
	//	TIME    sql.NullString
	//	ADD1    sql.NullString
	//	ADD2    sql.NullString
}

type ReData struct {
	C int64
	M string
	D interface{}
}
type CbData struct {
	Code int64
	Msg  string
	Data []int64
}
//订单系统传来的数据
type G_ORDERS_DATA struct {
	UsrId    int64
	Status   int64
	GoodsIds []int64
	Token    string
}
//加入购物车或直接购买传来的数据
type G_ADD_CART struct {
	GoodsId  int64
	GoodsNum int64
	Status   int64
	//status=20时为直接购买
}

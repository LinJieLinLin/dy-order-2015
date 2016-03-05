package common
import (
)


const (
	M1 = "数据库连接出错。"
	M2 = "读取用户信息失败。"
	M3 = "插入数据失败。"
	M4 = "读取数据失败。"
	M5 = "传入数据为空。"
	M6 = "传入数据有误。"
	M7 = "商品数据不存在，请刷新重试！"
	M8 = "无数据更新。"
	M9 = "页码超出范围"
	M10 = "暂无数据"
)

type ReData struct {
	C int64			`json:"code"`
	M string		`json:"msg"`
	D interface{}	`json:"data"`
}
//table Data
type O_CART struct {
	ID          int64
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
	TIME        string
	ADD1        string
	//	ADD2        string
}
//table orders
type O_ORDERS struct {
	ID          int64
	ORD_NO      string
	PUB_ORD_NO  string
	USR_ID      int64
	SHOP_ID     int64
	SHOP_NAME   string
	SHOP_URL    string
	SELLER_ID   int64
	SELLER_NAME string
	ORD_TYPE    int64
	STATUS      int64
	FROM_URL    string
	PRICE       float64
	CREATE_TIME string
	PAY_TIME    string
	TIME        string
	SYNC        int64
	//	ADD1        string
	//	ADD2        string
}
//table orders_item
type O_ORDERS_ITEM struct {
	ID          int64
	ORD_ID      int64
	GOODS_PRICE float64
	GOODS_ID    int64
	GOODS_NUM   int64
	GOODS_IMG   string
	GOODS_NAME  string
	GOODS_URL   string
	STATUS      int64
	TIME        string
	ADD1        string
	//	ADD2        string
}
//return cartData with FROM_URL
type RE_CART_DATA struct {
	FROM_URL string `json:"fromUrl"`
	SELLERS  []*SELLERS `json:"sellers"`
}
//return cartData
type SELLERS struct {
	FROM_URL    string `json:"fromUrl"`
	SELLER_ID   int64 `json:"sellerId"`
	SELLER_NAME string `json:"sellerName"`
	SHOP_ID     int64 `json:"shopId"`
	SHOP_NAME   string `json:"shopName"`
	SHOP_URL    string `json:"shopUrl"`
	GOODS       []*GOODS `json:"goods"`
}
type GOODS struct {
	ID          int64   `json:"id"`
	USER_ID     int64 `json:"userId"`
	SELLER_ID   int64 `json:"sellerId"`
	GOODS_ID    int64 `json:"goodsId"`
	GOODS_PRICE float64 `json:"goodsPrice"`
	GOODS_NUM   int64 `json:"goodsNum"`
	GOODS_IMG   string `json:"goodsImg"`
	GOODS_NAME  string `json:"goodsName"`
	GOODS_URL   string `json:"goodsUrl"`
	TIME        string `json:"time"`
}
//return ordersData
type RE_ORDERS_DATA struct {
	Page       int64 `json:"page"`
	OrdersList []*ORDERS `json:"ordersList"`
}
type ORDERS struct {
	ORD_ID      int64 `json:"ordId"`
	ORD_NO      string `json:"ordersNo"`
	SHOP_ID     int64 `json:"shopId"`
	SHOP_NAME   string `json:"shopName"`
	SHOP_URL    string `json:"shopUrl"`
	SELLER_ID   int64 `json:"sellerId"`
	SELLER_NAME string `json:"sellerName"`
	ORD_TYPE    string `json:"ordersType"`
	STATUS      int64 `json:"status"`
	FROM_URL    string `json:"fromUrl"`
	PRICE       float64 `json:"price"`
	CREATE_TIME string `json:"createTime"`
	GOODS       []*GOODS `json:"goods"`
}

type RE_CART struct {
	ID          int64   `json:"id"`
	USER_ID     int64 `json:"userId"`
	GOODS_ID    int64 `json:"goodsId"`
	GOODS_PRICE float64 `json:"goodsPrice"`
	GOODS_NUM   int64 `json:"goodsNum"`
	GOODS_IMG   string `json:"goodsImg"`
	GOODS_NAME  string `json:"goodsName"`
	GOODS_URL   string `json:"goodsUrl"`
	SHOP_ID     int64 `json:"shopId"`
	SHOP_NAME   string `json:"shopName"`
	SHOP_URL    string `json:"shopUrl"`
	SELLER_ID   int64 `json:"sellerId"`
	SELLER_NAME string `json:"sellerName"`
	FROM_URL    string `json:"fromUrl"`
	STATUS      int64 `json:"status"`
	TIME        string `json:"time"`
}
//GetCart data
type G_CartList struct {
	CartId  []int64
	FromUrl string
}
//edit cart Data
type G_CartData struct {
	CartId int64
	Num    int64
	Status int64
}
//submit order Data
type G_SubmitData struct {
	CartId []int64
}
//get orders data
type G_OrdersData struct {
	FromUrl string
	UserId  int64
	PNum    int64
	PSize   int64
	Status  int64
}
//return ordersData
type RE_ORDERS struct {
	OrderNo string `json:"orderNo"`
}
type ORDERS_ITEM struct {
	OrderId  int64
	FromUrl  string
	UsrId    int64
	Status   int64
	GoodsIds []int64
}
type ORDERS_CALLBACK_DATA struct {
	UsrId    int64
	Status   int64
	GoodsIds []int64
}
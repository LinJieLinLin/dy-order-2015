package alipay
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"databaseConn"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)
var db *sql.DB
var R *rand.Rand
const (
	del = 0
//	del = 1
)
func init() {
	databaseConn.CloseDB()
	var err error
	db, err= databaseConn.NewTestConn()
	if nil!=err {
		fmt.Println("连接数据库出错！！！！！！")
	}
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
}
func addData() {
	//o_orders
	temCount := 0
	if err := db.QueryRow("SELECT COUNT(ID) from o_orders  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败o_orders！", err)
		return
	}
	if temCount!=0 {
		fmt.Println("数据已存在。")
		return
	}
	if err := db.QueryRow("SELECT COUNT(ID) from o_orders_item  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败o_orders_item！", err)
		return
	}
	if temCount!=0 {
		fmt.Println("数据已存在。")
		return
	}
	addSql := `INSERT INTO o_orders (ID, ORD_NO,PUB_ORD_NO, USR_ID, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, ORD_TYPE, STATUS, FROM_URL, PRICE, CREATE_TIME, PAY_TIME, TIME, SYNC, ADD1, ADD2) VALUES
	('1', 'test','test', '1', '1', '测试店铺', '1', '1', '测试店主', '1', '0', 'http://orders.tmp.jxzy.com/api/EditO_USR?token=test', '1.00', '2015-07-17 15:24:26', '2015-07-17 15:24:29', '2015-07-17 15:24:32', '0', NULL, NULL);
`
	_, err := db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！1", err)
		return
	}
	addSql = `INSERT INTO o_orders_item (ID, ORD_ID, GOODS_PRICE, GOODS_ID, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, STATUS, TIME, ADD1, ADD2) VALUES
	('1', '1', '1.00', '1', '1', 'www.baidu.com', 'test', '1', '0', '2015-07-29 16:40:20', NULL, NULL);
`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！2", err)
		return
	}
	addSql = `INSERT INTO o_orders (ID, ORD_NO,PUB_ORD_NO, USR_ID, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, ORD_TYPE, STATUS, FROM_URL, PRICE, CREATE_TIME, PAY_TIME, TIME, SYNC, ADD1, ADD2) VALUES
	('2', 'test','test', '1', '1', '测试店铺', '1', '1', '测试店主', '1', '70', 'www.baidu.com', '0.00', '2015-07-17 15:24:26', '2015-07-17 15:24:29', '2015-07-17 15:24:32', '0', NULL, NULL);
`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！3", err)
		return
	}
	addSql = `INSERT INTO o_orders_item (ID, ORD_ID, GOODS_PRICE, GOODS_ID, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, STATUS, TIME, ADD1, ADD2) VALUES
	('2', '1', '0.00', '1', '1', 'www.baidu.com', 'test', '1', '0', '2015-07-29 16:40:20', NULL, NULL);
`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！4", err)
		return
	}
	//o_cart
	temCount=0
	if err := db.QueryRow("SELECT COUNT(ID) from o_cart  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败！",err)
		return
	}
	if temCount!=0{
		return
	}
	addSql = `
	INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('1', '36', '1', '10.00', '1', 'www.img.com', '测试课程', 'www.course.com', '1', '测试店铺', 'www.shop.com', '1', '测试卖家', 'http://orders.tmp.jxzy.com/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！",err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('2', '36', '2', '10.00', '1', 'www.img1.com', '测试课程1', 'www.course.com1', '2', '测试店铺1', 'www.shop.com1', '2', '测试卖家2', 'http://orders.tmp.jxzy.com/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！",err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('3', '36', '3', '10.00', '1', 'www.img2.com', '测试课程2', 'www.course.com1', '3', '测试店铺2', 'www.shop.com2', '3', '测试卖家3', 'http://orders.tmp.jxzy.com/api/getUrl', '20', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！",err)
	}
	fmt.Println("添加数据成功。")
}
func delData() int{
	re:=0
	if del!=0 {
		return re
	}
	delSql := "delete from o_orders where ORD_NO='test' or add1='test'"
	_, err := db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！1",err)
		return 1
	}
	delSql = "delete from o_orders_item where GOODS_NAME='test' or add1='test'"
	_, err = db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！2",err)
		return 1
	}
	delSql = "delete from o_cart where FROM_URL='FROM_URL' or add1='test'"
	_, err = db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！3",err)
		return 1
	}
	return re
}

func TestDel(t *testing.T){
	Convey("delData", t, func() {
		Convey("del", func() {
			re:=delData()
			So(re, ShouldEqual, 0)
		})
	})
}
func TestAdd(t *testing.T){
	Convey("addData", t, func() {
		Convey("add", func() {
			addData()
			So(nil, ShouldEqual, nil)
		})
	})
}
func TestGetOrdersPrice(t *testing.T) {
	Convey("GetOrdersPrice ", t, func() {
		Convey("success1", func() {
			ordNo:="test"
			price,e:=GetOrdersPrice(ordNo)
			fmt.Println("总价：",price)
			So(e, ShouldEqual,"")
		})
		Convey("e 没有该订单", func() {
			ordNo:="10086"
			price,e:=GetOrdersPrice(ordNo)
			fmt.Println("总价：",price)
			So(e, ShouldNotEqual,"")
		})
	})

}
func TestEditOrders(t *testing.T) {
	Convey("EditOrders",t,func(){
		Convey("success EditOrders", func() {
			ordNo:="test"
			e:=EditOrders(ordNo)
			So(e, ShouldEqual,"无数据更新。")
		})
		Convey("err1",func(){
			ordNo:="10086"
			e:=EditOrders(ordNo)
			fmt.Println(e)
			So(e, ShouldNotEqual,"")
		})
	})

}

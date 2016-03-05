package orders
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"common"
	"databaseConn"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
	"config"
)
var db *sql.DB
var R *rand.Rand
const (
	del = 0
//	del = 1
	delAll = 0
//	delAll = 1
)
func init() {
	config.TestInit("")
	config.ShowConf()
	databaseConn.CloseDB()
	var err error
	db, err= databaseConn.NewTestConn()
	if nil!=err {
		fmt.Println("连接数据库出错！！！！！！")
	}
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
}
func addData() {
	//o_cart
	temCount := 0
	if err := db.QueryRow("SELECT COUNT(ID) from o_cart  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败！", err)
		return
	}
	if temCount!=0 {
		return
	}
	addSql := `
	INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('1', '36', '1', '10.00', '1', 'www.img.com', '测试课程', 'www.course.com', '1', '测试店铺', 'www.shop.com', '1', '测试卖家', 'http://localhost:8001/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);`
	_, err := db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('2', '36', '2', '10.00', '1', 'www.img1.com', '测试课程1', 'www.course.com1', '2', '测试店铺1', 'www.shop.com1', '2', '测试卖家2', 'http://localhost:8001/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('3', '36', '3', '10.00', '1', 'www.img2.com', '测试课程2', 'www.course.com1', '3', '测试店铺2', 'www.shop.com2', '3', '测试卖家3', 'http://localhost:8001/api/getUrl', '20', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('4', '36', '3', '0.00', '1', 'www.img3.com', '测试课程3', 'www.course.com1', '4', '测试店铺3', 'www.shop.com3', '4', '测试卖家4', 'http://localhost:8001/api/getUrl', '20', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	//o_orders
	temCount = 0
	if err := db.QueryRow("SELECT COUNT(ID) from o_orders  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败o_orders！", err)
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
	addSql = `INSERT INTO o_orders (ID, ORD_NO,PUB_ORD_NO, USR_ID, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, ORD_TYPE, STATUS, FROM_URL, PRICE, CREATE_TIME, PAY_TIME, TIME, SYNC, ADD1, ADD2) VALUES
	('1', 'test','ptest', '1', '1', '测试店铺', '1', '1', '测试店主', '1', '0', 'http://orders.tmp.jxzy.com', '1.00', '2015-07-17 15:24:26', '2015-07-17 15:24:29', '2015-07-17 15:24:32', '0', NULL, NULL);
`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
		return
	}
	addSql = `INSERT INTO o_orders_item (ID, ORD_ID, GOODS_PRICE, GOODS_ID, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, STATUS, TIME, ADD1, ADD2)
	VALUES ('1', '1', '1.00', '1', '1', 'www.baidu.com', 'test', '1', '1', '2015-07-29 16:40:20', NULL, NULL);
`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
		return
	}
	fmt.Println("添加数据成功。")
}
func delData() int64 {
	re := int64(0)
	if del!=0 {
		return re
	}
	delSql := "delete from o_orders "
	if delAll==0 {
		delSql+="where (ORD_NO='test' or add1='test' )"
	}
	_, err := db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！1", err)
		return 1
	}
	delSql = "delete from o_orders_item  "
	if delAll==0 {
		delSql+=" where GOODS_NAME='test' or add1='test'"
	}
	_, err = db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！2", err)
		return 1
	}
	delSql = "delete from o_cart "
	if delAll==0 {
		delSql+=" where FROM_URL='FROM_URL' or add1='test'"
	}
	_, err = db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！", err)
		return 1
	}
	fmt.Println("删除成功！")
	return re
}

func TestDel(t *testing.T) {
	Convey("delData", t, func() {
		Convey("del", func() {
			re := delData()
			So(re, ShouldEqual, 0)
		})
	})
}
func TestAdd(t *testing.T) {
	Convey("addData", t, func() {
		Convey("add", func() {
			addData()
			So(nil, ShouldEqual, nil)
		})
	})
}
func TestSubmitOrderF(t *testing.T) {
	Convey("SubmitOrderF ", t, func() {
		Convey("success1", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{1}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			fmt.Println(re.D)
			So(0, ShouldEqual, re.C)
		})
		Convey("success2直接购买", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{3}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			fmt.Println(re.D)
			So(0, ShouldEqual, re.C)
		})
		Convey("success3直接购买", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{4}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			fmt.Println(re.D)
			So(0, ShouldEqual, re.C)
		})
		Convey("不存在订单号", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{10086}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			So(1, ShouldEqual, re.C)
		})
		Convey("用户ID有误", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{1}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, -1)
			So(1, ShouldEqual, re.C)
		})
		Convey("err数据有误1", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{[]int64{-2}}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			So(1, ShouldEqual, re.C)
		})
		Convey("err数据有误2", func() {
			arg_d := []*common.G_SubmitData{}
			re := SubmitOrderF(arg_d, 36)
			So(1, ShouldEqual, re.C)
		})
		Convey("err数据有误3", func() {
			arg_d := []*common.G_SubmitData{}
			d := common.G_SubmitData{}
			arg_d = append(arg_d, &d)
			re := SubmitOrderF(arg_d, 36)
			So(1, ShouldEqual, re.C)
		})
	})
}
func TestGetOrdersDataF(t *testing.T) {
	Convey("GetOrdersDataF ", t, func() {
		Convey("success0", func() {
			arg_d := common.G_OrdersData{}
			arg_d.UserId = 36
			arg_d.FromUrl = ""
			arg_d.PNum = 1
			arg_d.PSize = 1
			arg_d.Status = 0
			re := GetOrdersDataF(arg_d)
			fmt.Println(re.D)
			So(0, ShouldEqual, re.C)
		})
		Convey("success1", func() {
			arg_d := common.G_OrdersData{}
			arg_d.UserId = 1
			arg_d.FromUrl = "http://orders"
			arg_d.PNum = 1
			arg_d.PSize = 2
			arg_d.Status = 0
			re := GetOrdersDataF(arg_d)
			fmt.Println(re.D)
			So(0, ShouldEqual, re.C)
		})
		Convey("err用户ID", func() {
			arg_d := common.G_OrdersData{}
			arg_d.UserId = -1
			arg_d.FromUrl = "http://orders"
			arg_d.PNum = 1
			arg_d.PSize = 1
			arg_d.Status = 0
			re := GetOrdersDataF(arg_d)
			fmt.Println(re.D)
			So(1, ShouldEqual, re.C)
		})
		Convey("err每页数量", func() {
			arg_d := common.G_OrdersData{}
			arg_d.UserId = 1
			arg_d.FromUrl = "http://orders"
			arg_d.PNum = 1
//			arg_d.PSize = 1
			arg_d.Status = 0
			re := GetOrdersDataF(arg_d)
			fmt.Println(re.D)
			So(1, ShouldEqual, re.C)
		})
		Convey("err页数", func() {
			arg_d := common.G_OrdersData{}
			arg_d.UserId = 1
			arg_d.FromUrl = "http://orders"
//			arg_d.PNum = -1
			arg_d.PSize = 1
			arg_d.Status = 0
			re := GetOrdersDataF(arg_d)
			fmt.Println(re.D)
			So(1, ShouldEqual, re.C)
		})

	})
}
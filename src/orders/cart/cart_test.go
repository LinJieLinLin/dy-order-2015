package cart
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
	"github.com/Centny/gwf/routing"
	"net/http"
	"net/url"
	"github.com/Centny/gwf/util"
)
var db *sql.DB
var R *rand.Rand
const (
	del = 0
//	del = 1
)
var mux *routing.SessionMux
func init() {
	config.TestInit("")
	config.ShowConf()
	databaseConn.CloseDB()
	var err error
	db, err= databaseConn.NewTestConn()
	if nil!=err {
		fmt.Println("连接数据库出错！！！！！！")
		return
	}
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
	//set test server mux
	sb := routing.NewSrvSessionBuilder("", "/", "orders", 6000000, 1000)
	mux = routing.NewSessionMux("", sb)
	mux.HFilterFunc("^/.*$", common.TestSetSession)
	mux.HFunc("^/l/addToCart(\\?.*)?$", AddToCart)
	mux.HFunc("^/l/getCart(\\?.*)?$", GetCart)
	mux.HFunc("^/l/editCart(\\?.*)?$", EditCart)
}
func addData() {
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
						('1', '36', '1', '10.00', '1', 'www.img.com', '测试课程', 'www.course.com', '1', '测试店铺', 'www.shop.com', '1', '测试卖家', 'http://orders.tmp.jxzy.com/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);`
	_, err := db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('2', '36', '2', '10.00', '1', 'www.img1.com', '测试课程1', 'www.course.com1', '2', '测试店铺1', 'www.shop.com1', '2', '测试卖家2', 'http://orders.tmp.jxzy.com/api/getUrl', '0', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('3', '36', '3', '10.00', '1', 'www.img2.com', '测试课程2', 'www.course.com1', '3', '测试店铺2', 'www.shop.com2', '3', '测试卖家3', 'http://orders.tmp.jxzy.com/api/getUrl', '20', '2015-07-15 16:25:02', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
		INSERT INTO o_cart (ID, USER_ID, GOODS_ID, GOODS_PRICE, GOODS_NUM, GOODS_IMG, GOODS_NAME, GOODS_URL, SHOP_ID, SHOP_NAME, SHOP_URL, SELLER_ID, SELLER_NAME, FROM_URL, STATUS, TIME, ADD1, ADD2) VALUES
						('4', '70', '1106100792', '135.00', '1', 'http://localhost:8082/resImg/1106100380/1106100490_1437732752092239600_imgPre0.jpg', '女装', 'http://localhost:8082/resCenter.html#/Res/1106100490', '1106100380', '公主家', 'http://localhost:8082/decorate/storeFirstPage.html?1106100380', '1111061384', '123123', 'http://localhost:8082/api/getUrl', '20', '2015-07-28 10:22:39', 'test', NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
}
func delData() {
	if del!=0 {
		return
	}
	delSql := "delete from o_cart where FROM_URL='FROM_URL' or add1='test'"
	_, err := db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！")
	}
}
func TestDel(t *testing.T) {
	Convey("delData", t, func() {
		Convey("del", func() {
			delData()
			So(nil, ShouldEqual, nil)
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

func TestAddToCartF(t *testing.T) {
	Convey("AddToCartF", t, func() {
		Convey("success ", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 36, 1,1 , 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 10086, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(0, ShouldEqual, re.C)
		})
		Convey("success update", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 36, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 10086, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(0, ShouldEqual, re.C)
		})

		Convey("success 直接购买", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 36, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 10086, "SELLER_NAME", "FROM_URL", 20, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(0, ShouldEqual, re.C)
		})
		Convey("error 商品状态有误 ", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 36, 1,1 , 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 10086, "SELLER_NAME", "FROM_URL", 2, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("没有传入数据", func() {
			cartDates := []*common.O_CART{}
//			cartData := common.O_CART{1, -1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
//			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("用户ID有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, -1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品ID有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, -1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品价格有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, -1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品数量有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 0, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品名称有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "", "", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品链接有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "", 1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("商店ID有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", -1, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("卖家ID有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 0, "SHOP_NAME", "SHOP_URL", 0, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("卖家有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 0, "SHOP_NAME", "SHOP_URL", 1, "", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
		Convey("来源URL有误", func() {
			cartDates := []*common.O_CART{}
			cartData := common.O_CART{1, 1, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 0, "SHOP_NAME", "SHOP_URL", 1, "SELLER_NAME", "", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			re := AddToCartF(cartDates)
			So(1, ShouldEqual, re.C)
		})
	})
}
func TestGetCartF(t *testing.T) {
	Convey("GetCartF ", t, func() {
		Convey("success 读取全部数据", func() {
			data:=common.G_CartList{}
			re := GetCartF(36,data)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("success 包含URL 有数据", func() {
			data:=common.G_CartList{}
			data.FromUrl="http://"
			re := GetCartF(36,data)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("success 包含URL 无数据", func() {
			data:=common.G_CartList{}
			data.FromUrl="-1"
			re := GetCartF(36,data)
			fmt.Println("返回数据：",re)
			So("购物车为空！", ShouldEqual, re.M)
		})
		Convey("success 获取直接购买的商品", func() {
			data:=common.G_CartList{}
			data.CartId = []int64{3,2,1}
			re := GetCartF(36,data)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("error 商品ID有误", func() {
			data:=common.G_CartList{}
			data.CartId = []int64{-1,2,1}
			re := GetCartF(36,data)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("error 用户ID", func() {
			data:=common.G_CartList{}
			data.CartId = []int64{3,2,1}
			re := GetCartF(0,data)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
	})
}
func TestEditCartF(t *testing.T) {
	Convey("EditCartF ", t, func() {
		Convey("su", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,2,0}
			arg_d = append(arg_d,&d)

			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("su 有事务,修改状态改为30", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,1,30}
			arg_d = append(arg_d,&d)
			db,err:= databaseConn.GetConn()
			if err!=nil{
				So(1, ShouldEqual, -1)
			}
			tx,_:=db.Begin()
			re := EditCartF(arg_d,36,tx)
			tx.Commit()
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("su 修改状态改为0", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,1,-3}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("su 修改状态改为30", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{2,1,30}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
		Convey("error 状态30状态改为0", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,1,-3}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})

		Convey("购物车无数据", func() {
			arg_d := []*common.G_CartData{}
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("购物车ID有误", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{0,1,0}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("商品数量有误", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,-1,0}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("用户ID有误", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,1,0}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,0,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("购物车状态有误", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{1,1,2}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(1, ShouldEqual, re.C)
		})
		Convey("su 删除购物车", func() {
			arg_d := []*common.G_CartData{}
			d := common.G_CartData{2,1,-1}
			arg_d = append(arg_d,&d)
			re := EditCartF(arg_d,36,nil)
			fmt.Println("返回数据：",re)
			So(0, ShouldEqual, re.C)
		})
	})
}

func TestHttpFunc(t *testing.T) {
	Convey("TestHttpFunc", t, func() {
		Convey("TestHttpFunc", func() {
			c := []*http.Cookie{}
			c = append(c, &http.Cookie{Name:"uid", Value: "36", })
			c = append(c, &http.Cookie{Name:"username", Value: "linjie", })
			re := common.ReData{}

			v := url.Values{}
			cartDates := []*common.O_CART{}
			data := ""
			cartData := common.O_CART{1, 36, 1, 1, 1, "GOODS_IMG", "GOODS_NAME", "GOODS_URL", 1, "SHOP_NAME", "SHOP_URL", 10086, "SELLER_NAME", "FROM_URL", 0, "2015-07-15 15:42:45", "test"}
			cartDates = append(cartDates, &cartData)
			data=util.S2Json(cartDates)
			v.Add("data", data)

			//addToCart
			re=common.TestRequestWithHs(mux, "POST", "/l/addToCart", v, c)
			So(re.C, ShouldEqual, 0)
			re=common.TestRequestWithHs(mux, "POST", "/l/addToCart", v, nil)
			So(re.C, ShouldEqual, 1)
			re=common.TestRequestWithHs(mux, "POST", "/l/addToCart", nil, c)
			So(re.C, ShouldEqual, 1)
			//getCart
			v.Set("data", "{}")
			re=common.TestRequestWithHs(mux, "GET", "/l/getCart", v, c)
			So(re.C, ShouldEqual, 0)
			v.Set("data", "{from:-1}")
			re=common.TestRequestWithHs(mux, "GET", "/l/getCart", v, c)
			So(re.C, ShouldEqual, 1)
			//editCart
			editCart := []*common.G_CartData{}
			d := common.G_CartData{1, 1, 30}
			editCart = append(editCart, &d)
			data=util.S2Json(editCart)
			v.Set("data", data)
			re=common.TestRequestWithHs(mux, "GET", "/l/editCart", v, c)
			So(re.C, ShouldEqual, 0)
			v.Set("data", "{}")
			re=common.TestRequestWithHs(mux, "GET", "/l/editCart", v, c)
			So(re.C, ShouldEqual, 1)
		})
	})
}

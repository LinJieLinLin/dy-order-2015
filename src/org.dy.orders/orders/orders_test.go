package orders
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
)
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
}
func addData() {
	//o_cart
	temCount := 0
	if err := db.QueryRow("SELECT COUNT(ID) from o_usr  WHERE ID=1").Scan(&temCount); err != nil {
		fmt.Println("查找数据失败！", err)
		return
	}
	if temCount!=0 {
		return
	}
	addSql := `
	INSERT INTO o_usr (ID, USR_ID, GOODS_ID, STATUS, TIME, TOKEN, ADD1, ADD2) VALUES ('1', '36', '1', '0', '2015-07-23 12:07:00', 'test', NULL, NULL);
	`
	_, err := db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	addSql = `
	INSERT INTO o_usr (ID, USR_ID, GOODS_ID, STATUS, TIME, TOKEN, ADD1, ADD2) VALUES ('2', '36', '2', '0', '2015-07-23 12:07:00', 'test', NULL, NULL);
	`
	_, err = db.Exec(addSql)
	if nil!=err {
		fmt.Println("添加数据失败！", err)
	}
	fmt.Println("添加数据成功。")
}
func delData() int64 {
	re := int64(0)
	if del!=0 {
		return re
	}
	fmt.Printf("test db con--:%v",config.TestDbConfig())

	if nil==db{
		fmt.Println("连接数据库出错！！！！！！")
	}
	delSql := "delete from o_usr where (TOKEN='test' or add1='test' )"
	_, err := db.Exec(delSql)
	if nil!=err {
		fmt.Println("删除数据失败！1", err)
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
func TestGetUrlF(t *testing.T) {
	Convey("GetUrlF ", t, func() {
		Convey("GetUrlF su", func() {
			arg_d := G_ORDERS_DATA{36,0,[]int64{1},"test"}
			token,re := GetUrlF(db,arg_d)
			fmt.Println(re,token)
			So("", ShouldEqual, re)
		})
		Convey("传入数据有误商品ID", func() {
			arg_d := G_ORDERS_DATA{36,0,[]int64{0},"test"}
			token,re := GetUrlF(db,arg_d)
			fmt.Println(re,token)
			So("", ShouldNotEqual, re)
		})
		Convey("传入数据有误人员ID", func() {
			arg_d := G_ORDERS_DATA{0,0,[]int64{1},"test"}
			token,re := GetUrlF(db,arg_d)
			fmt.Println(re,token)
			So("", ShouldNotEqual, re)
		})

	})
}
//func TestSyncOrdersF(t *testing.T) {
//	Convey("SyncOrdersF ", t, func() {
//		Convey("success1", func() {
//			arg_d := G_ORDERS_DATA{36,10,[]int64{1},"test"}
//			re := SyncOrdersF(db,arg_d)
//			fmt.Println(re)
//			So("", ShouldEqual, re)
//		})
//	})
//}
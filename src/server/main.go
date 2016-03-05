package server

import (
	"config"
	"databaseConn"
	"log"
	"net/http"
	"fmt"
	"os"
	"bufio"
	"github.com/Centny/gwf/util"
	cny4goLog "github.com/Centny/gwf/log"
	"org.dy.orders/orders"
	"io"
	"sync"
	C "common"
)

var onceApp sync.Once

func startOnce() {
	go func() {
		//port
		config.ShowConf()
		srvPort := config.ServerPort()
		//dbConnect
		dbConfig := config.DbConfig()
		//log
		logPath := config.LogPath()
		logName := config.LogFileName()
		//create and open log file
		fmt.Println("log path : "+ logPath+logName)
		util.FTouch(logPath+logName)
		f, err := os.OpenFile(logPath+logName, os.O_RDWR  | os.O_CREATE  |os.O_APPEND, 0666)
		if err != nil {
			panic(err.Error())
			return
		}
		//create file buffer writer.
		bo := bufio.NewWriter(f)
		defer f.Close()
		defer bo.Flush()
		cny4goLog.SetWriter(io.MultiWriter(bo, os.Stdout))
		//set log level（1 show all）
		cny4goLog.SetLevel(1)
		//log end

		databaseConn.SetConnConfig(dbConfig)
		db, err := databaseConn.GetNewConn()
		if nil!=err{
			cny4goLog.E("connet db failed--:%v",err)
			return
		}
		err=C.CheckT(db)
		if nil!=err{
			cny4goLog.E("auto create ordersTable failed--:%v",err)
			return
		}
		err=orders.CheckT(db)
		if nil!=err{
			cny4goLog.E("auto create o_usr failed--:%v",err)
			return
		}
		RouteConfig()
		if err := http.ListenAndServe(":"+srvPort, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func startTestOnce() {
	go func() {
		RouteConfig()
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func Start() {
	onceApp.Do(startOnce)
}

func StartTestServer() {
	onceApp.Do(startTestOnce)
}

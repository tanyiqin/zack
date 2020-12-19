package main

import (
	"flag"
	"fmt"
	mdb "github.com/tanyiqin/zack/db"
	"github.com/tanyiqin/zack/router"
	"github.com/tanyiqin/zack/znet"
	"os"
	"os/signal"
	"runtime"
)

func StopFunc(znet.IConnection) {
	fmt.Println("stopFunc ....")
}

var numCores = flag.Int("n", 2, "number of CPU cores to use")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*numCores)

	s := znet.NewServer()

	// 添加Router
	s.AddRouter(1, router.CsAccountLogin)
	s.AddRouter(2, router.CsPlayerCreate)

	// 添加链接中断前执行的操作 一般为保存玩家数据
	s.AddConnStopFunc(StopFunc)

	// 服务器启动
	go s.Start()

	// 这里捕获退出信号 执行需要的退出操作
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c

	fmt.Println("server stop with sig=", sig)
	mdb.Mdb.StopMongoDB()
	s.GetConnMgr().Stop()
}

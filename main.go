package main

import (
	"flag"
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/conf"
	"github.com/labstack/echo/v4"
	"runtime"
	"shy2you/api"
)

func main() {
	runtime.GOMAXPROCS(conf.ServerConf.MaxProc)
	flag.Parse()
	// 关闭redis
	defer ironman.Redis.Close()
	e := echo.New()
	api.Routers(e)
	ironman.Server(e)
}

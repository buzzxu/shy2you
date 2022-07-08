package api

import (
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/conf"
	"github.com/buzzxu/ironman/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"shy2you/api/ws"
	"shy2you/internal/say"
	"shy2you/pkg/types"
)

func init() {
	//
	conf.LoadDefaultConf()
	logger.InitLogger()
	ironman.RedisConnect()
	ironman.DefaultJWTConfig.AuthScheme = conf.ServerConf.Jwt.AuthScheme
	ironman.DefaultJWTConfig.ContextKey = conf.ServerConf.Jwt.ContextKey
	ironman.DefaultJWTConfig.SigningMethod = conf.ServerConf.Jwt.SigningMethod
	ironman.DefaultJWTConfig.TokenLookup = "query:" + echo.HeaderAuthorization
	ironman.DefaultJWTConfig.Claims = &types.Claims{}
	go say.Start()

}
func Routers(e *echo.Echo) {
	ironman.JwtConfig(middleware.DefaultSkipper)
	e.GET("/notify/ws", ws.Notify)
	e.POST("/notify/say", ws.Say)
	e.POST("/notify/ping", ws.Ping)
	//e.Use(middleware.JWTWithConfig(jwtConfig))
}

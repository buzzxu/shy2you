package ws

import (
	"context"
	boystypes "github.com/buzzxu/boys/types"
	"github.com/buzzxu/ironman/jwtt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"net/http"
	"shy2you/pkg/auth"
	"shy2you/pkg/types"
	"shy2you/pkg/websockets"
	"strconv"
	"time"
)

var (
	SessionsPool = websockets.New()
	upgrader     = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//r.Response.Header.Add("Access-Control-Allow-Origin", "*")
			//r.Response.Header.Add("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD, PUT, POST")
			return true
		},
	}
	corsHeaders = http.Header{}
)

type PingSendMessage struct {
	UserId  int64  `json:"userId"`
	Context string `json:"context"`
}

func init() {
	corsHeaders.Add("Access-Control-Allow-Origin", "*")
	corsHeaders.Add("Access-Control-Allow-Credentials", "true")
	corsHeaders.Add("Access-Control-Expose-Headers", "Content-Disposition, Authorization, Content-Type, x-requested-with, GET, POST, OPTIONS, PUT, DELETE")
}
func Ping(c echo.Context) error {
	var msg PingSendMessage
	if err := c.Bind(&msg); err != nil {
		return c.JSON(200, boystypes.ErrorOf(err))
	}
	userId := strconv.FormatInt(msg.UserId, 10)
	err := SessionsPool.Send(userId, "发射成功:"+msg.Context)
	if err != nil {
		return c.JSON(500, boystypes.ErrorOf(err))
	}
	return c.JSON(200, boystypes.ResultOf(200, true))
}

func Say(c echo.Context) error {
	var say types.Say
	if err := c.Bind(&say); err != nil {
		return c.JSON(200, boystypes.ErrorOf(err))
	}
	c.Logger().Infof(" receive message userId: %s", say.UserId)
	err := SessionsPool.Say(&say)
	if err != nil {
		return c.JSON(500, boystypes.ErrorOf(err))
	}
	return c.JSON(200, boystypes.ResultOf(200, true))
}

func Notify(c echo.Context) error {
	token, err := jwtt.ParserTokenUnverified(c, echo.HeaderAuthorization)
	if err != nil {
		return err
	}
	if token.IsPresent() {
		claims := token.Get().(*jwt.Token).Claims.(*types.Claims)
		//find user info
		var session *websockets.Session
		session, err = auth.GetUserSession(claims)
		if err == nil {
			ws, err := upgrader.Upgrade(c.Response(), c.Request(), corsHeaders)
			if err != nil {
				return err
			}
			SessionsPool.Lock()
			SessionsPool.Sessions[ws] = *session
			defer func(connection *websocket.Conn) {
				SessionsPool.Lock()
				delete(SessionsPool.Sessions, connection)
				SessionsPool.Unlock()
			}(ws)
			SessionsPool.Unlock()
			defer func(ws *websocket.Conn) {
				err := ws.Close()
				if err != nil {
					c.Logger().Error(err)
				}
			}(ws)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
							return
						}
					}
				}
			}()
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					c.Logger().Error(err)
					ws.Close()
					break
				}
			}
		}

	}

	return nil
}

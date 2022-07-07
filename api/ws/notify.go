package ws

import (
	"fmt"
	boystypes "github.com/buzzxu/boys/types"
	"github.com/buzzxu/ironman"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"net/http"
	"shy2you/pkg/auth"
	"shy2you/pkg/types"
	"shy2you/pkg/websockets"
	"strconv"
)

var (
	SessionsPool = websockets.New()
	upgrader     = websocket.Upgrader{}
)

type PingSendMessage struct {
	UserId  int64  `json:"userId"`
	Context string `json:"context"`
}

func Ping(c echo.Context) error {
	var msg PingSendMessage
	if err := c.Bind(&msg); err != nil {
		return c.JSON(200, boystypes.ErrorOf(err))
	}
	userId := strconv.FormatInt(msg.UserId, 10)
	SessionsPool.Send(userId, "发射成功:"+msg.Context)
	return c.JSON(200, boystypes.ResultOf(200, true))
}

func Notify(c echo.Context) error {
	token, err := ironman.ParserTokenUnverified(c, echo.HeaderAuthorization, "")
	if err != nil {
		return err
	}
	if token.IsPresent() {
		claims := token.Get().(*jwt.Token).Claims.(*types.Claims)
		print(claims.Subject)
		//find user info
		var session *websockets.Session
		session, err = auth.GetUserSession(claims)
		if err == nil {
			upgrader.CheckOrigin = func(r *http.Request) bool {
				return true
			}
			ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
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
			defer ws.Close()
			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					c.Logger().Error(err)
				}
				fmt.Println(msg)
			}
		}

	}

	return nil
}

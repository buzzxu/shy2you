package ws

import (
	"fmt"
	boystypes "github.com/buzzxu/boys/types"
	"github.com/buzzxu/ironman"
	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"shy2you/pkg/auth"
	"shy2you/pkg/types"
	"shy2you/pkg/websockets"
	"strconv"
)

var SessionsPool = websockets.New()

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
			websocket.Handler(func(ws *websocket.Conn) {
				SessionsPool.Lock()
				SessionsPool.Sessions[ws] = *session
				defer func(connection *websocket.Conn) {
					SessionsPool.Lock()
					delete(SessionsPool.Sessions, connection)
					SessionsPool.Unlock()
				}(ws)
				SessionsPool.Unlock()
				defer ws.Close()
				msg := ""
				for {
					if err := websocket.Message.Receive(ws, &msg); err != nil {
						c.Logger().Error(err)
						return
					}
					fmt.Println(msg)
				}
			}).ServeHTTP(c.Response(), c.Request())
		}

	}

	return nil
}

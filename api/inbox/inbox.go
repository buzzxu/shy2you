package inbox

import (
	"fmt"
	boystypes "github.com/buzzxu/boys/types"
	"github.com/buzzxu/ironman"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"net/http"
	"shy2you/pkg/auth"
	"shy2you/pkg/inbox"
	"shy2you/pkg/types"
	"strconv"
)

var (
	SessionsPool = inbox.New()
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
	SessionsPool.Send(userId, "发射成功:"+msg.Context)
	return c.JSON(200, boystypes.ResultOf(200, true))
}

func Dispatch(c echo.Context) error {
	var message types.InboxMessage
	if err := c.Bind(&message); err != nil {
		return c.JSON(200, boystypes.ErrorOf(err))
	}
	c.Logger().Infof(" receive message userId: %s", message.UserId)
	SessionsPool.Dispatch(&message)
	return c.JSON(200, boystypes.ResultOf(200, true))
}

func Notify(c echo.Context) error {
	token, err := ironman.ParserTokenUnverified(c, echo.HeaderAuthorization, "")
	if err != nil {
		return err
	}
	if token.IsPresent() {
		claims := token.Get().(*jwt.Token).Claims.(*types.Claims)
		//find user info
		var session *inbox.Session
		session, err = auth.GetInboxUser(claims)
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

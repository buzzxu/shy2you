package inbox

import (
	"encoding/json"
	"github.com/buzzxu/ironman/logger"
	"github.com/gorilla/websocket"
	"shy2you/pkg/types"
	"sync"
)

type InboxDispatcher struct {
	sync.RWMutex
	Sessions map[*websocket.Conn]Session
}
type Session struct {
	UserId string
}

func New() *InboxDispatcher {
	return &InboxDispatcher{
		Sessions: make(map[*websocket.Conn]Session),
	}
}

func (s *InboxDispatcher) Send(userId, message string) error {
	s.RLock()
	defer s.RUnlock()
	for con, session := range s.Sessions {
		if session.UserId == userId {
			if err := con.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (s *InboxDispatcher) Dispatch(drop *types.InboxDrop) error {
	s.RLock()
	defer s.RUnlock()
	data, err := json.Marshal(drop.Data)
	if err != nil {
		return err
	}
	if data == nil {
		logger.Of("ws").Info("say nothing.")
		return nil
	}
	message := string(data[:])
	for con, session := range s.Sessions {
		// UserId not null
		if drop.UserId != "" && session.UserId == drop.UserId {
			logger.Infof("user: %s ,message: %s", session.UserId, drop.Data)
			if err := con.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

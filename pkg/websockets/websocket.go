package websockets

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"shy2you/pkg/types"
	"sync"
)

type SessionPool struct {
	sync.RWMutex
	Sessions map[*websocket.Conn]Session
}
type Session struct {
	UserId     string
	Type       int
	CompanyId  int
	SupplierId int
	TenantId   int
}

func New() *SessionPool {
	return &SessionPool{
		Sessions: make(map[*websocket.Conn]Session),
	}
}

func (s *SessionPool) Send(userId, message string) error {
	s.RLock()
	defer s.RUnlock()
	for con, session := range s.Sessions {
		if session.UserId == userId {
			if err := websocket.Message.Send(con, message); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (s *SessionPool) Say(say *types.Say) error {
	s.RLock()
	defer s.RUnlock()
	data, err := json.Marshal(say.Data)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	message := string(data[:])
	for con, session := range s.Sessions {
		// UserId not null
		if say.UserId != "" && session.UserId == say.UserId {
			if err := websocket.Message.Send(con, message); err != nil {
				return err
			}
			break
		} else {
			switch say.Region {
			case 1:
				//tenantId
				if session.TenantId == say.CompanyId {
					if err := websocket.Message.Send(con, message); err != nil {
						return err
					}
				}
				break
			case 2:
				//companyId
				if session.CompanyId == say.CompanyId {
					if err := websocket.Message.Send(con, message); err != nil {
						return err
					}
				}
				break
			case 3:
				//supplierId
				if session.SupplierId == say.CompanyId {
					if err := websocket.Message.Send(con, message); err != nil {
						return err
					}
				}
				break
			}
		}
	}
	return nil
}
func (s *SessionPool) SendAll(message string) error {
	s.RLock()
	defer s.RUnlock()
	for session := range s.Sessions {
		if err := websocket.Message.Send(session, message); err != nil {
			return err
		}
	}
	return nil
}
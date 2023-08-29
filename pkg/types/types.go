package types

import "github.com/dgrijalva/jwt-go"

type (
	// JWT 信息
	Claims struct {
		Type     int    `json:"type"`
		UserName string `json:"userName"`
		Region   int    `json:"region"`
		jwt.StandardClaims
	}
	Say struct {
		UserId    string      `json:"userId"`
		Region    int         `json:"region"`
		Types     []int       `json:"types"`
		Type      int         `json:"type"`
		CompanyId int         `json:"companyId"`
		Data      interface{} `json:"data"`
	}

	InboxDrop struct {
		UserId string          `json:"userId"`
		Data   []*InboxMessage `json:"data"`
	}

	InboxMessage struct {
		Id        string      `json:"id""`
		UserId    int         `json:"userId"`
		Status    int         `json:"status"`
		ObjId     string      `json:"objId"`
		Region    string      `json:"region"`
		BizType   string      `json:"bizType"`
		Title     string      `json:"title"`
		Content   string      `json:"content"`
		Path      string      `json:"path"`
		Data      interface{} `json:"data"`
		Time      string      `json:"time"`
		CreatedAt string      `json:"createdAt"`
		UpdatedAt string      `json:"updatedAt"`
	}
)

func (s *Say) IsRegion(t int) bool {
	if s.Types != nil {
		for _, val := range s.Types {
			if val == t {
				return true
			}
		}
	}
	return false
}

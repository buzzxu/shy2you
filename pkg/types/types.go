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

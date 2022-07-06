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
		CompanyId int         `json:"companyId"`
		Data      interface{} `json:"data"`
	}
)

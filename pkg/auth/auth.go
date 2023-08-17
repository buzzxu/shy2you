package auth

import (
	"context"
	"errors"
	"github.com/buzzxu/ironman"
	"shy2you/pkg/inbox"
	"shy2you/pkg/types"
	"shy2you/pkg/websockets"
	"strconv"
)

func GetUserSession(claims *types.Claims) (*websockets.Session, error) {
	var opt = GetUser(claims)
	if opt.IsPresent() {
		hash := opt.Get().(map[string]string)
		userType, _ := strconv.Atoi(hash["type"])
		var tenantId, companyId, supplierId int
		if hash["tenantId"] != "" {
			tenantId, _ = strconv.Atoi(hash["tenantId"])
		}
		if hash["companyId"] != "" {
			companyId, _ = strconv.Atoi(hash["companyId"])
		}
		if hash["supplierId"] != "" {
			supplierId, _ = strconv.Atoi(hash["supplierId"])
		}
		return &websockets.Session{
			UserId:     hash["id"],
			Type:       userType,
			CompanyId:  companyId,
			SupplierId: supplierId,
			TenantId:   tenantId,
		}, nil
	}
	return nil, errors.New("not found user in redis")
}

func GetInboxUser(claims *types.Claims) (*inbox.Session, error) {
	var opt = GetUser(claims)
	if opt.IsPresent() {
		hash := opt.Get().(map[string]string)
		return &inbox.Session{
			UserId: hash["id"],
		}, nil
	}
	return nil, errors.New("not found user in redis")
}

func GetUser(claims *types.Claims) *ironman.Optional {
	var ctx = context.Background()
	hash := ironman.Redis.HGetAll(ctx, "user:"+claims.Subject)
	val, err := hash.Result()
	if err != nil {
		return ironman.OptionalOfNil()
	}
	return ironman.OptionalOf(val)
}

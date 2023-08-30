package commons

import (
	"context"
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/logger"
)

func CreateStreamExists(ctx context.Context, topic, group string) {
	exits, err := ironman.Redis.Exists(ctx, topic).Result()
	if err != nil {
		return
	}
	if exits == 0 {
		_, err := ironman.Redis.XGroupCreateMkStream(context.Background(), topic, group, "$").Result()
		if err != nil {
			logger.Infof("create stream fail, %s", err.Error())
			return
		}
	} else {
		logger.Infof("stream topic: %s,group:%s,OK", topic, group)
	}
}

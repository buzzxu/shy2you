package handlers

import (
	"context"
	"encoding/json"
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/logger"
	"github.com/redis/go-redis/v9"
	"shy2you/api/ws"
	"shy2you/pkg/commons"
	"shy2you/pkg/types"
)

func Say() {
	var topic = "topic:shy2you:notify"
	var group = "shy2you"
	commons.CreateStreamExists(topic, group)
	for {
		var ctx = context.Background()
		logger.Infof("Ready receive new message")
		datas, err := ironman.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: "say",
			Streams:  []string{topic, ">"},
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			logger.Errorf("receive new message error. %s", err.Error())
			continue
		}
		logger.Infof("receive new message")
		for _, result := range datas {
			for _, message := range result.Messages {
				messageID := message.ID
				values := message.Values
				var say = types.Say{}
				err := json.Unmarshal([]byte(values["data"].(string)), &say)
				if err != nil {
					logger.Errorf("parser message error: %s , del it", err.Error())
					ironman.Redis.XDel(ctx, topic, group, messageID)
					continue
				}
				ws.SessionsPool.Say(&say)
				ironman.Redis.XAck(ctx, topic, group, messageID)
			}
		}
		//for i := 0; i < len(datas[0].Messages); i++ {
		//	messageID := datas[0].Messages[i].ID
		//	values := datas[0].Messages[i].Values
		//	var say = types.Say{}
		//	json.Unmarshal([]byte(values["data"].(string)), &say)
		//	ws.SessionsPool.Say(&say)
		//	ironman.Redis.XAck(ctx, topic, group, messageID)
		//}
	}

}

package handlers

import (
	"context"
	"encoding/json"
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/logger"
	"github.com/redis/go-redis/v9"
	"shy2you/api/inbox"
	"shy2you/pkg/commons"
	"shy2you/pkg/types"
)

func Inbox() {
	var topic = "topic:shy2you:inbox"
	var group = "shy2you"
	commons.CreateStreamExists(topic, group)
	for {
		var ctx = context.Background()
		logger.Infof("Ready receive new inbox message")
		datas, err := ironman.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: "inbox",
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
				var inboxDrop = types.InboxDrop{}
				err := json.Unmarshal([]byte(values["data"].(string)), &inboxDrop)
				if err != nil {
					logger.Errorf("parser message error: %s , del it", err.Error())
					ironman.Redis.XDel(ctx, topic, group, messageID)
					continue
				}
				err = inbox.SessionsPool.Dispatch(&inboxDrop)
				if err != nil {
					logger.Errorf("message send error: %s , del it", err.Error())
					return
				}
				ironman.Redis.XAck(ctx, topic, group, messageID)
			}
		}
	}

}

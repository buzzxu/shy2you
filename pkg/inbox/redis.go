package inbox

import (
	"context"
	"encoding/json"
	"github.com/buzzxu/ironman"
	"github.com/buzzxu/ironman/logger"
	"github.com/redis/go-redis/v9"
	"shy2you/pkg/types"
	"strconv"
	"time"
)

// 获取最新的消息
func FetchLatestUnRead(userId string, send func(inboxDrop *types.InboxDrop)) {
	ctx := context.Background()
	messageIds, _ := ironman.Redis.ZRevRangeByScore(ctx, "inbox:messages:"+userId, &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	var messages []*types.InboxMessage
	cmds := make([]*redis.MapStringStringCmd, len(messageIds))
	pipeline := ironman.Redis.Pipeline()
	for i, messageId := range messageIds {
		cmds[i] = pipeline.HGetAll(ctx, "inbox:message:"+messageId)
	}
	_, err := pipeline.Exec(ctx)
	if err != nil {
		logger.Errorf("批量获取数据失败原因: {}", err.Error())
		return
	}

	for _, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			logger.Errorf("获取收件箱消息失败,原因: {}", err.Error())
			continue
		}
		message := convertInboxMessage(data)
		message.UserId, _ = strconv.Atoi(userId)
		messages = append(messages, message)
	}
	if len(messages) > 0 {
		send(&types.InboxDrop{UserId: userId, Data: messages})
	}
	//
	cleanExpireMessage(ctx, userId)
}

func convertInboxMessage(data map[string]string) *types.InboxMessage {
	var message types.InboxMessage
	var err error
	message.Id = data["id"]
	message.Status, err = strconv.Atoi(data["status"])
	if err != nil {
		message.Status = 0
	}
	message.ObjId = data["objId"]
	message.Region = data["region"]
	message.BizType = data["bizType"]
	message.Title = data["title"]
	message.Content = data["content"]
	message.Path = data["path"]

	var _data interface{}
	err = json.Unmarshal([]byte(data["data"]), &_data)
	if err != nil {
		message.Data = nil
	}
	message.Data = _data
	message.Time = data["time"]
	message.CreatedAt = data["createdAt"]
	message.UpdatedAt = data["updatedAt"]
	return &message
}

func cleanExpireMessage(ctx context.Context, userId string) {
	// 获取当前时间
	now := time.Now()
	// 计算过期时间的时间戳
	expireTimestamp := now.Add(-7 * 24 * time.Hour).Unix()
	messageIds, err := ironman.Redis.ZRangeByScore(ctx, "inbox:messages:"+userId, &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(expireTimestamp, 10)}).Result()
	if err != nil {
		logger.Errorf("清理过期消息失败,原因: %s", err.Error())
		return
	}
	if len(messageIds) == 0 {
		return
	}
	key_messages := "inbox:messages:" + userId
	cmds := make([]*redis.IntCmd, len(messageIds))
	delCmds := make([]*redis.IntCmd, len(messageIds))
	pipeline := ironman.Redis.Pipeline()
	pipelineDel := ironman.Redis.Pipeline()
	for i, messageId := range messageIds {
		cmds[i] = pipeline.ZRem(ctx, key_messages, messageId)
		delCmds[i] = pipeline.Del(ctx, "inbox:message:"+messageId)
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return
	}
	_, err = pipelineDel.Exec(ctx)
	if err != nil {
		return
	}
}

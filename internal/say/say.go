package say

import (
	"context"
	"encoding/json"
	"github.com/buzzxu/ironman"
	"github.com/go-redis/redis/v8"
	"shy2you/api/ws"
	"shy2you/pkg/types"
)

var topic = "topic:shy2you:notify"
var group = "jedi"

func init() {

}
func Start() {
	statusCmd := ironman.Redis.XGroupCreateMkStream(context.Background(), topic, group, "$")
	if statusCmd.Err() != nil {
		return
	}
	for {
		var ctx = context.Background()
		datas, err := ironman.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: "say",
			Streams:  []string{topic, ">"},
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			break
		}
		for _, result := range datas {
			for _, message := range result.Messages {
				messageID := message.ID
				values := message.Values
				var say = types.Say{}
				json.Unmarshal([]byte(values["data"].(string)), &say)
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

package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

const (
	Room       = "room:global"
	HistoryKey = "history:global"
)

func Init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatal("Erro ao conectar no Redis: ", err)
	}
}

func Publish(msg Message) {

	data, _ := json.Marshal(msg)
	Rdb.Publish(Ctx, Room, data)
}

func SaveHistory(msg Message) {

	data, _ := json.Marshal(msg)

	pipe := Rdb.TxPipeline()
	pipe.RPush(Ctx, HistoryKey, data)
	pipe.LTrim(Ctx, HistoryKey, -50, -1)
	pipe.Expire(Ctx, HistoryKey, 1*time.Minute)
	pipe.Exec(Ctx)
}

func GetHistory() []Message {
	var messages []Message

	items, err := Rdb.LRange(Ctx, HistoryKey, 0, -1).Result()
	if err != nil {
		return messages
	}

	for _, item := range items {
		var msg Message
		json.Unmarshal([]byte(item), &msg)
		messages = append(messages, msg)
	}

	return messages
}

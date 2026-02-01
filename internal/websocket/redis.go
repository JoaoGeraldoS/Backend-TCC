package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		opt = &redis.Options{
			Addr: redisUrl,
		}
	}

	Rdb = redis.NewClient(opt)

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
	now := time.Now().Unix()
	expirationLimit := 60

	pipe := Rdb.TxPipeline()

	pipe.ZAdd(Ctx, HistoryKey, redis.Z{
		Score:  float64(now),
		Member: data,
	})

	pipe.ZRemRangeByScore(Ctx, HistoryKey, "-inf", fmt.Sprintf("%d", expirationLimit))
	pipe.ZRemRangeByRank(Ctx, HistoryKey, 0, -51)
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

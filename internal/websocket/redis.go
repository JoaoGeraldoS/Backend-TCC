package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

const (
	Room       = "room:global"
	HistoryKey = "history:global:v4"
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

func Publish(msg map[string]interface{}) {

	data, _ := json.Marshal(msg)
	Rdb.Publish(Ctx, Room, data)
}

func SaveHistory(msg Message) string {
	uniqueID := uuid.New().String()

	msgData := map[string]interface{}{
		"user":    msg.User,
		"message": msg.Message,
		"time":    msg.Time,
		"id":      uniqueID,
	}

	data, _ := json.Marshal(msgData)

	now := time.Now().Unix()

	err := Rdb.ZAdd(Ctx, HistoryKey, redis.Z{
		Score:  float64(now),
		Member: data,
	}).Err()

	if err != nil {
		log.Println("Erro ao salvar no Redis:", err)
	}

	limiteTempo := now - 120
	Rdb.ZRemRangeByScore(Ctx, HistoryKey, "-inf", fmt.Sprintf("%d", limiteTempo))

	log.Printf("Mensagem salva. Limpeza de mensagens anteriores a: %d", limiteTempo)

	return uniqueID
}

func GetHistory() []Message {
	var messages []Message

	now := time.Now().Unix()

	Rdb.ZRemRangeByScore(Ctx, HistoryKey, "-inf", fmt.Sprintf("%d", now-120))

	items, err := Rdb.ZRange(Ctx, HistoryKey, 0, -1).Result()
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

package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return rdb
}

func PublishMessage(ctx context.Context, rdb *redis.Client, room, msg string) error {
	return rdb.Publish(ctx, room, msg).Err()
}

func SubscribeRoom(ctx context.Context, rdb *redis.Client, room string) {
	sub := rdb.Subscribe(ctx, room)
	ch := sub.Channel()

	for msg := range ch {
		println("Messagen recebida:", msg.Payload)
	}
}

func SaveMessage(ctx context.Context, rdb *redis.Client, room string, msg string) error {
	key := "history:" + room

	pipe := rdb.TxPipeline()
	pipe.RPush(ctx, key, msg)
	pipe.LTrim(ctx, key, -50, -1) // últimas 50
	pipe.Expire(ctx, key, 24*time.Hour)
	_, err := pipe.Exec(ctx)

	return err
}

func UserJoin(ctx context.Context, rdb *redis.Client, room, user string) {
	key := "online:" + room
	rdb.SAdd(ctx, key, user)
	rdb.Expire(ctx, key, 1*time.Hour)
}

func UserLeave(ctx context.Context, rdb *redis.Client, room, user string) {
	key := "online:" + room
	rdb.SRem(ctx, key, user)
}

func GetOnline(ctx context.Context, rdb *redis.Client, room string) ([]string, error) {
	return rdb.SMembers(ctx, "online:"+room).Result()
}

func GetHistory(ctx context.Context, rdb *redis.Client, room string) ([]string, error) {
	key := "history:" + room
	return rdb.LRange(ctx, key, 0, -1).Result()
}

func mains() {
	var ctx = context.Background()

	rdb := NewRedis(ctx)
	PublishMessage(ctx, rdb, "room:Golang", "Joao: Ola!")
	go SubscribeRoom(ctx, rdb, "room:Golang")
}

package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedis(uri string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Redis is connected at %s \n", uri)

	return client
}

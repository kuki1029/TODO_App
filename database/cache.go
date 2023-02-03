package database

import (
	"github.com/go-redis/redis"

	"fmt"
)

// Sets up the redis cache
func RedisSetUp() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return client
}

// Checks if the server is online for redis
func Ping(client *redis.Client) error {
	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

package middleware

import (
	"strings"
	"time"

	"github.com/go-redis/redis"

	"fmt"
	"strconv"
)

// RedisClient is a repository for interacting with the redis cache
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new instance of the Redis Client
func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisClient{client: client}
}

// Checks if the server is online for redis
func (rc RedisClient) Ping() error {
	pong, err := rc.client.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

// This function gets the users details from redis storage
func (rc RedisClient) GetFromRedis(key string) (uint, error) {
	userDetails := rc.client.Get(key)
	user, err := userDetails.String(), userDetails.Err()
	if (err != nil) || (len(string(user)) == 0) {
		// 0 does not represent any id in the database so we return that incase of an error.
		return 0, err
	}

	idString := strings.Fields(user)
	userID, _ := strconv.ParseUint(idString[5], 10, 64)
	ID := uint(userID)
	return ID, nil
}

// Sets the value in the redis cache
func (rc RedisClient) SetInRedis(value string, redisVal string, timeAmt time.Duration) {
	rc.client.Set(value, redisVal, timeAmt)
}

// Deletes the value in the redis cache
func (rc RedisClient) DelInRedis(value string) {
	rc.client.Del(value)
}

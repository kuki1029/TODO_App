package database

import (
	"strings"

	"github.com/go-redis/redis"

	"fmt"
	"strconv"
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

// This function gets the users details from redis storage
func GetFromRedis(client *redis.Client, key string) (uint, error) {
	userDetails := client.Get(key)
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

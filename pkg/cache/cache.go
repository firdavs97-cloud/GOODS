package cache

import (
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var rc *Cache

// RedisCache represents a Redis cache client
type Cache struct {
	client *redis.Client
}

// Connect creates a new instance of RedisCache
func Connect(addr, password string) {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // Redis server address
		Password: password,
	})

	// Ping Redis to ensure connection is established
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatalln(err)
	}

	rc = &Cache{client: rdb}
}

// Set sets a value in the cache with an expiration time
func Set(key string, value interface{}, expiration time.Duration) error {
	err := rc.client.Set(rc.client.Context(), key, value, expiration).Err()
	log.Println(err)
	return err
}

// Get retrieves a value from the cache
func Get(key string) (string, error) {
	return rc.client.Get(rc.client.Context(), key).Result()
}

// FlushAll clears all keys from the cache
func FlushAll() error {
	return rc.client.FlushAll(rc.client.Context()).Err()
}

// Delete deletes a key from the cache
func Delete(key string) error {
	return rc.client.Del(rc.client.Context(), key).Err()
}

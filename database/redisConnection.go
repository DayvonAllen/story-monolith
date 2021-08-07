package database

import (
	"github.com/gomodule/redigo/redis"
)

var Conn *redis.Pool

func ConnectToRedis() {
	const maxConnections = 10
	redisPool := &redis.Pool{
		MaxIdle: maxConnections,
		Dial:    func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") },
	}
	Conn = redisPool
}

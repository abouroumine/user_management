package redis

import "github.com/gomodule/redigo/redis"

// Redis ...
var Redis redis.Conn

// InitRedis ...
func InitRedis() {
	con, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	Redis = con
}

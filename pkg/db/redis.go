package db

import "github.com/xochat/xochat_im_server/pkg/redis"

func NewRedis(addr string, password string) *redis.Conn {
	return redis.New(addr, password)
}

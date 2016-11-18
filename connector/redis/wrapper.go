package rrredis

import (
	"./redis"
	"errors"
	"fmt"
	"rrframework/logs"
	"time"
)

type RedisOptions struct {
	dialTimeout  time.Duration
	db           int
	password     string
	connPoolSize int
}

type clientPool struct {
}

// connect once when program start
func Connect(addr string, opt *RedisOptions) error {
	if addr == "" || opt == nil {
		return fmt.Errorf("Redis addr empty or options nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:        opt.addr,
		DB:          opt.db,
		DialTimeout: opt.dialTimeout,
		PoolSize:    connPoolSize,
	})
	if client == nil {
		return fmt.Errorf(fmt.Sprintf("Connect to redis [%s] fail", addr))
	}
}

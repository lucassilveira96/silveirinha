package database

import (
	"context"
	"fmt"
	"sync"
	"teste/internal/infra/logger"
	"teste/internal/infra/logger/attributes"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Redis struct {
	rdb     *redis.Client
	opt     *redis.Options
	retries []time.Duration
	locker  sync.Mutex
}

func NewRedis(c *fiber.Ctx, opt *redis.Options, lazyConnection bool) *Redis {
	rdb := &Redis{
		opt: opt,
		retries: []time.Duration{
			250 * time.Millisecond,
			500 * time.Millisecond,
			1000 * time.Millisecond,
			2500 * time.Millisecond,
			5000 * time.Millisecond,
		},
	}

	if !lazyConnection {
		_ = rdb.initializeAndGetRedis(c)
	}

	return rdb
}

func (r *Redis) Connection(c *fiber.Ctx) *redis.Client {
	return r.initializeAndGetRedis(c)
}

func (r *Redis) Close(c *fiber.Ctx) {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.rdb == nil {
		return
	}

	if err := r.rdb.Close(); err != nil {
		logger.Error(
			c,
			"Failed to close Redis",
			r.configToAttribute().WithError(err),
		)
	}

	r.rdb = nil
}

func (r *Redis) initializeAndGetRedis(c *fiber.Ctx) *redis.Client {
	rdb := r.rdb
	if rdb != nil {
		return rdb
	}

	r.locker.Lock()
	defer r.locker.Unlock()

	// double-checked locking
	if rdb = r.rdb; rdb != nil {
		return rdb
	}

	start := time.Now()
	logger.Info(
		c,
		"Initializing Redis",
		r.configToAttribute(),
	)

	rdb = redis.NewClient(r.opt)
	connected := false
	var err error

	for retry, duration := range r.retries {
		if err = r.checkConnection(rdb); err != nil {
			logger.Warn(
				c,
				fmt.Sprintf("Connection retry [%d]: Redis connection", retry+1),
				r.configToAttribute().WithError(err),
			)
			time.Sleep(duration)
		} else {
			connected = true
			err = nil
			break
		}
	}

	if !connected {
		if err = r.checkConnection(rdb); err != nil {
			logger.Fatal(
				c,
				"Failed to connect to Redis database",
				r.configToAttribute().WithError(err),
			)
		}
	}

	elapsed := time.Since(start)
	logger.Info(
		c,
		fmt.Sprintf("Redis initialized in [%v]", elapsed),
		r.configToAttribute(),
	)
	r.rdb = rdb
	return rdb
}

func (r *Redis) checkConnection(rdb *redis.Client) error {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rdb.Ping(timeout).Err()
}

func (r *Redis) configToAttribute() attributes.Attributes {
	config := r.opt
	return attributes.Attributes{
		"redis.address":  config.Addr,
		"redis.db":       config.DB,
		"redis.password": "[Masked]",
	}
}

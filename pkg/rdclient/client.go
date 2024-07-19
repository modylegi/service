package rdclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	numAttempts   = 5
	retryInterval = 5 * time.Second
)

func NewClient(ctx context.Context, cfg *Config) (*redis.Client, error) {
	var count int
	connStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s", cfg.username, cfg.password, cfg.host, cfg.port, cfg.database)
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		for count = 0; count < numAttempts; count++ {
			client = redis.NewClient(opt)
			if err := client.Ping(ctx).Err(); err != nil {
				fmt.Printf("can't connect to cache attempt: %d error: %s\n", count+1, err.Error())
				time.Sleep(retryInterval)
			} else {
				return client, nil
			}
		}
		return nil, fmt.Errorf("exceeded the number of connection attempts connstr=%s", connStr)

	}
	return client, nil
}

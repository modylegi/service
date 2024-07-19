package pgconn

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	numAttempts   = 5
	retryInterval = 5 * time.Second
)

func NewConnection(ctx context.Context, cfg *Config) (*sqlx.DB, error) {
	var count int
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.username, cfg.password, cfg.host, cfg.port, cfg.database)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		for count = 0; count < numAttempts; count++ {
			db, err = sqlx.Open("pgx", connStr)
			if err != nil {
				return nil, err
			}
			if err := db.PingContext(ctx); err != nil {
				fmt.Printf("can't connect to base attempt: %d error: %s\n", count+1, err.Error())
				time.Sleep(retryInterval)
			} else {
				return db, nil
			}
		}
		return nil, fmt.Errorf("exceeded the number of connection attempts connstr=%s", connStr)
	}
	return db, nil
}

package database

import (
	"fmt"
	"sync"
	"time"

	//revive:disable-next-line
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

var (
	db   *sqlx.DB
	once sync.Once
)

func NewPostgresPool(host, port, user, pass, dbName, sslMode string) *sqlx.DB {
	once.Do(func() {
		dataSourceName := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, pass, dbName, sslMode,
		)

		pool, err := sqlx.Connect("pgx", dataSourceName)
		if err != nil {
			log.Fatal(nil, map[string]interface{}{
				"error": err,
			}, "failed to connect to database")
		}

		pool.SetMaxOpenConns(100)
		pool.SetMaxIdleConns(10)
		pool.SetConnMaxLifetime(60 * time.Minute)

		db = pool
	})

	return db
}

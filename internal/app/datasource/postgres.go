package datasource

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/otter-im/profile/internal/app/config"
	"log"
	"os"
	"sync"
	"time"
)

var (
	dbPool   *sql.DB
	poolOnce sync.Once
)

func DB() *sql.DB {
	poolOnce.Do(func() {
		db, err := sql.Open("postgres", config.Config().DatabaseURL)
		if err != nil {
			log.Printf("db establish failure: %v", err)
			os.Exit(2)
		}
		dbPool = db
	})
	return dbPool
}

func CheckDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := DB().PingContext(ctx); err != nil {
		return fmt.Errorf("db establish failure: %v", err)
	}
	return nil
}

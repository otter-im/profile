package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/otter-im/profile/internal/config"
	"sync"
	"time"
)

var (
	pgMainOnce sync.Once
	pgMain     *pg.DB
)

func Postgres() *pg.DB {
	pgMainOnce.Do(func() {
		options := &pg.Options{
			Addr:     config.PostgresAddress(),
			User:     config.PostgresUser(),
			Password: config.PostgresPassword(),
			Database: config.PostgresDatabase(),
		}

		if config.PostgresSSL() {
			options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		pgMain = pg.Connect(options)
		AddExitHook(func() error {
			if err := pgMain.Close(); err != nil {
				return err
			}
			return nil
		})
	})
	return pgMain
}

func checkPostgres() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := Postgres().Ping(ctx); err != nil {
		return fmt.Errorf("postgresql connection failure: %v", err)
	}
	return nil
}

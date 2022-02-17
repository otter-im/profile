package main

import (
	"flag"
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/otter-im/profile/internal/app"
	"github.com/otter-im/profile/internal/app/datasource"
	"log"
	"os"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

var (
	exitHooks = make([]func() error, 0)
)

func main() {
	flag.Usage = usage
	flag.Parse()

	err := app.Init()
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	oldVersion, newVersion, err := migrations.Run(datasource.Postgres(), flag.Args()...)
	if err != nil {
		log.Printf("migration %d -> %d failed: %s\n", oldVersion, newVersion, err)
		os.Exit(-3)
	}

	if newVersion != oldVersion {
		fmt.Printf("migrated from %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", newVersion)
	}

	err = app.Exit()
	if err != nil {
		log.Print(err)
		os.Exit(-2)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

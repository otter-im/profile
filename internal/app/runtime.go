package app

import (
	"github.com/gorilla/mux"
	"github.com/otter-im/profile/internal/config"
	"golang.org/x/exp/rand"
	"log"
	mathRand "math/rand"
	"net"
	"net/http"
	"time"
)

var (
	exitHooks = make([]func() error, 0)
)

func Init() error {
	log.Printf("environment: %s\n", config.ServiceEnvironment())

	rand.Seed(uint64(time.Now().UnixNano()))
	mathRand.Seed(time.Now().UnixNano())

	if err := checkPostgres(); err != nil {
		return err
	}

	if err := checkRedis(); err != nil {
		return err
	}
	return nil
}

func Run() error {
	router := mux.NewRouter()

	http.Handle("/", router)
	addr := net.JoinHostPort(config.ServiceHost(), config.ServicePort())
	log.Printf("listening on: %v\n", addr)
	return http.ListenAndServe(addr, nil)
}

func Exit() error {
	for _, hook := range exitHooks {
		err := hook()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddExitHook(hook func() error) {
	exitHooks = append(exitHooks, hook)
}

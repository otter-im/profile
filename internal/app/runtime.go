package app

import (
	"github.com/gorilla/mux"
	"github.com/otter-im/profile/internal/app/config"
	"github.com/otter-im/profile/internal/app/datasource"
	"github.com/otter-im/profile/internal/app/handler"
	"golang.org/x/exp/rand"
	"log"
	mathRand "math/rand"
	"net/http"
	"os"
	"time"
)

func Init() error {
	log.Printf("environment: %s\n", config.Config().ServiceEnv)

	rand.Seed(uint64(time.Now().UnixNano()))
	mathRand.Seed(time.Now().UnixNano())

	if err := datasource.CheckDB(); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	if err := datasource.CheckRedis(); err != nil {
		log.Println(err)
		os.Exit(2)
	}
	return nil
}

func Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/profile", handler.ProfileDefaultHandler)
	router.HandleFunc("/profile/{user_id}", handler.ProfileByIdHandler)

	http.Handle("/", router)
	log.Printf("listening on: %v\n", config.Config().ServiceAddress)
	return http.ListenAndServe(config.Config().ServiceAddress, nil)
}

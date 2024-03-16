package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/phuwn/crawlie/src/db"
	"github.com/phuwn/crawlie/src/handler"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/service"
	"github.com/phuwn/crawlie/src/store"
	_ "gorm.io/driver/postgres"
)

func main() {
	config, err := configLoad()
	if err != nil {
		log.Fatal(err)
	}

	err = db.NewPostgresConn(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	store := store.New()
	service := service.New(config.Service)
	server.NewServerCfg(store, service)

	addr := fmt.Sprintf(":%d", config.Server.ListeningPort)
	log.Printf("listening on port%s\n", addr)

	err = http.ListenAndServe(addr, handler.Router(config.Server))
	if err != nil {
		log.Printf("server got terminated, err: %s\n", err.Error())
	}
}

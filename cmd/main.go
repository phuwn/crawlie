package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/phuwn/crawlie/src/config"
	"github.com/phuwn/crawlie/src/router"
	"github.com/phuwn/crawlie/src/server"
	_ "gorm.io/driver/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Router.ListeningPort)
	log.Printf("listening on port%s\n", addr)

	err = http.ListenAndServe(addr, router.Router())
	if err != nil {
		log.Printf("server got terminated, err: %s\n", err.Error())
	}
}

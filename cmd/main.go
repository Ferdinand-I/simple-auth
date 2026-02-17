package main

import (
	"authapp/internal/config"
	"authapp/internal/server"
	"authapp/internal/storage"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	db, err := storage.ConnectDB(cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
		server := server.New()
	server.SetUp(db)
	server.Run(":" + cfg.Server.Port)
}

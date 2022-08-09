package main

import (
	"log"

	"github.com/osalomon89/test-crud-api/internal/config"
	"github.com/osalomon89/test-crud-api/internal/repository"
	"github.com/osalomon89/test-crud-api/internal/server"
)

func main() {
	config.Load()
	dbConn, err := repository.GetConnectionDB()
	if err != nil {
		log.Fatal("Can't connect mysql db. ", err)
	}

	if err := repository.Migrate(dbConn); err != nil {
		log.Fatal("Can't create table on mysql db. ", err)
	}

	s := server.Server{
		DBConn: dbConn,
		Port:   config.Port(),
	}

	s.Start()
}

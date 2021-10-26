package main

import (
	"context"
	"final-project/internal/config"
	"final-project/internal/db/postgres"
	"final-project/internal/repository"
	"fmt"
	"log"
)

func main() {
	cfgFile, err := config.LoadConfig(".config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	psqlDB, err := postgres.NewPsqlDb(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", psqlDB.Stat())
	}
	defer psqlDB.Close()

	ctx := context.Background()
	url := "http://www.cbr.ru/scripts/XML_daily.asp?date_req=25/10/2021"
	fmt.Println(repository.ParseXml(url))
	parsed := repository.ParseXml(url)
	r := repository.Instance{
		Db: psqlDB,
	}
	err = r.Insert(ctx, parsed)
	if err != nil {
		log.Fatalf("Insert error: %#v", err)
	}
}

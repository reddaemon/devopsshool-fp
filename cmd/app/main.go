package main

import (
	"final-project/internal/config"
	"final-project/internal/db/postgres"
	"final-project/internal/repository"
	"final-project/internal/usecase"

	//"fmt"
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

	r := repository.NewInstance(psqlDB)
	uc := usecase.NewUseCase(r)

	//uc.PullDataByPeriod("28/10/2021")
	uc.GetData()
}

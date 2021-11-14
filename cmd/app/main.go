package main

import (
	"final-project/internal/config"
	"final-project/internal/db/postgres"
	"final-project/internal/db/redis"
	"final-project/internal/handlers"
	"final-project/internal/repository"
	"final-project/internal/router"
	"final-project/internal/usecase"
	"net/http"

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
	redisConn, err := redis.NewRedisConn(cfg)
	if err != nil {
		log.Fatalf("Redis connection init: %s", err)
	}

	r := repository.NewInstance(psqlDB, redisConn)
	uc := usecase.NewUseCase(r)

	h := handlers.NewHandler(uc)
	router := router.RegisterRouter(h)
	http.ListenAndServe(":8080", router)
}

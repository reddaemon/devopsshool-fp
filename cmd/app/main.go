package main

import (
	"database/sql"
	"final-project/internal/config"
	"final-project/internal/db/postgres"
	"final-project/internal/db/redis"
	"final-project/internal/handlers"
	_ "final-project/internal/migrations"
	"final-project/internal/repository"
	"final-project/internal/router"
	"final-project/internal/usecase"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"net/http"

	"log"
)

func main() {
	cfgFile, err := config.LoadConfig()
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

	mdb, _ := sql.Open("postgres", psqlDB.Config().ConnString())
	err = goose.Up(mdb, "./internal/migrations")
	if err != nil {
		log.Printf("migration err: %#v", err)
		log.Printf("%s", err)
	}

	redisConn, err := redis.NewRedisConn(cfg)
	if err != nil {
		log.Fatalf("Redis connection init: %s", err)
	}

	r := repository.NewInstance(psqlDB, redisConn)
	uc := usecase.NewUseCase(r)

	h := handlers.NewHandler(uc)
	router := router.RegisterRouter(h)
	http.ListenAndServe(":8085", router)

}

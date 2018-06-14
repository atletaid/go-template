package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tokopedia/sqlt"

	_ "github.com/lib/pq"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/account/delivery/rest"
	"github.com/sportivaid/go-template/src/account/repository/postgres"
	"github.com/sportivaid/go-template/src/account/repository/redis"
	"github.com/sportivaid/go-template/src/common/auth"

	"github.com/sportivaid/go-template/config"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ldate)

	//Init Config
	cfg, ok := config.InitConfig([]string{"files/etc/config"}...)
	if !ok {
		fmt.Println("Error opening config files")
		return
	}

	flag.Parse()

	// Init Database Postgres
	dbMaster, err := sqlt.Open("postgres", cfg.Account.MasterDB)
	if err != nil {
		log.Println("Error opening database : ", err)
		return
	}

	// Init Redis
	redisPool, err := redis.NewPool(cfg.Redis.Host, cfg.Redis.DialTimeout*time.Second, cfg.Redis.IdleTimeout*time.Second, cfg.Redis.PoolSize)
	if err != nil {
		log.Println(err)
		return
	}

	authUsecase := auth.NewUsecase()
	authMiddleware := auth.NewMiddleware(authUsecase)

	accountRepo := postgres.NewAccountRepository(dbMaster, dbMaster, cfg.Server.DBTimeout*time.Second)
	accountCache := redis.NewAccountCache(cfg.InMemory.DefaultExpiration, cfg.InMemory.IntervalPurges)
	accountRepo = redis.NewMiddlewareAccountRepository(accountCache, redisPool, accountRepo)
	accountUsecase := account.NewAccountUsecase(accountRepo)
	router := rest.NewUserHandler(authMiddleware, accountUsecase)
	router.Run(cfg.Account.Port)
}

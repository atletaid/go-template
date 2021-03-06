package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/atletaid/go-template/config"
	"github.com/atletaid/go-template/src/common/auth"
	"github.com/atletaid/go-template/src/module/account"
	"github.com/atletaid/go-template/src/module/account/delivery"
	"github.com/atletaid/go-template/src/module/account/repository"
	"github.com/tokopedia/sqlt"
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

	// Init PostgreSQL Database
	dbMaster, err := sqlt.Open("postgres", cfg.Account.MasterDB)
	if err != nil {
		log.Println("Error opening database : ", err)
		return
	}

	// Init Inmemory & Redis Cache
	accountCache := repository.NewAccountCache(cfg.InMemory.DefaultExpiration, cfg.InMemory.IntervalPurges)
	redisPool, err := repository.NewPool(cfg.Redis.Host, cfg.Redis.DialTimeout*time.Second, cfg.Redis.IdleTimeout*time.Second, cfg.Redis.PoolSize)
	if err != nil {
		log.Println(err)
		return
	}

	authMiddleware := auth.NewMiddleware()

	accountRepo := repository.NewAccountRepository(dbMaster, dbMaster, cfg.Server.DBTimeout*time.Second)
	accountRepo = repository.NewMiddlewareAccountRepository(accountCache, redisPool, accountRepo)
	accountUsecase := account.NewAccountUsecase(accountRepo)

	var ginRouter *gin.Engine
	if cfg.Server.Enviroment == "development" {
		ginRouter = gin.Default()
	} else {
		ginRouter = gin.New()
	}

	router := delivery.NewAccountHandler(ginRouter, authMiddleware, accountUsecase)
	router.Run(cfg.Account.Port)
}

package redis

import (
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func NewPool(host string, dialTimeout time.Duration, idleTimeout time.Duration, poolSize int) (*redigo.Pool, error) {
	pool := redigo.Pool{
		MaxActive:   poolSize,
		MaxIdle:     poolSize,
		IdleTimeout: idleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", host, redigo.DialConnectTimeout(dialTimeout))
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return c, err
		},
	}

	if _, err := pool.Dial(); err != nil {
		pool.Close()
		log.Println(err)
		return nil, err
	}

	return &pool, nil
}

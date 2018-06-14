package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
)

var (
	KeyAccountsFindAll = "accounts:findall"
	KeyAccountsFind    = "accounts:find"
)

type accountRepository struct {
	cache map[string]*cache.Cache
	pool  *redigo.Pool
	next  account.AccountRepository
}

func NewAccountCache(cExpiration time.Duration, cIntervalPurges time.Duration) map[string]*cache.Cache {
	return map[string]*cache.Cache{
		KeyAccountsFindAll: cache.New(cExpiration*time.Minute, cIntervalPurges*time.Minute),
		KeyAccountsFind:    cache.New(cExpiration*time.Minute, cIntervalPurges*time.Minute),
	}
}

func NewMiddlewareAccountRepository(cache map[string]*cache.Cache, pool *redigo.Pool, next account.AccountRepository) account.AccountRepository {
	return &accountRepository{
		cache: cache,
		pool:  pool,
		next:  next,
	}
}

func (ar *accountRepository) FindAll() ([]*model.Account, error) {
	field := "accounts"

	accountCache, found := ar.cache[KeyAccountsFindAll].Get(field)
	if found {
		return accountCache.([]*model.Account), nil
	}

	accountsJSON, err := redigo.Bytes(ar.do("HGET", KeyAccountsFindAll, field))

	if err == redigo.ErrNil && ar.next != nil {
		accounts, err := ar.next.FindAll()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		accountsJSON, err := json.Marshal(&accounts)
		if err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := ar.do("HSET", KeyAccountsFindAll, field, accountsJSON); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := ar.do("EXPIRE", KeyAccountsFindAll, 3600); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		ar.cache[KeyAccountsFindAll].SetDefault(field, accounts)
		return accounts, nil
	}

	accounts := make([]*model.Account, 0)
	if err := json.Unmarshal(accountsJSON, &accounts); err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	ar.cache[KeyAccountsFindAll].SetDefault(field, accounts)
	return accounts, nil
}

func (ar *accountRepository) Find(userID int64) (*model.Account, error) {
	field := fmt.Sprintf("%v", userID)

	accountCache, found := ar.cache[KeyAccountsFind].Get(field)
	if found {
		return accountCache.(*model.Account), nil
	}

	accountJSON, err := redigo.Bytes(ar.do("HGET", KeyAccountsFind, field))
	if err == redigo.ErrNil && ar.next != nil {
		account, err := ar.next.Find(userID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		accountJSON, err := json.Marshal(&account)
		if err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := ar.do("HSET", KeyAccountsFind, field, accountJSON); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := ar.do("EXPIRE", KeyAccountsFind, 3600); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		ar.cache[KeyAccountsFind].SetDefault(field, account)
		return account, nil
	}

	var account *model.Account
	if err := json.Unmarshal(accountJSON, &account); err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	ar.cache[KeyAccountsFind].SetDefault(field, account)
	return account, nil
}

func (ar *accountRepository) do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := ar.pool.Get()
	defer conn.Close()

	return conn.Do(command, args...)
}

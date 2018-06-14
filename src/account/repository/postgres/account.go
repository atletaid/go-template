package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
	"github.com/tokopedia/sqlt"
)

type accountRepository struct {
	DbMaster *sqlt.DB
	DbSlave  *sqlt.DB
	Timeout  time.Duration
}

func NewAccountRepository(dbMaster *sqlt.DB, dbSlave *sqlt.DB, timeout time.Duration) account.AccountRepository {
	return &accountRepository{
		DbMaster: dbMaster,
		DbSlave:  dbSlave,
		Timeout:  timeout,
	}
}

func (ar *accountRepository) FindAll() ([]*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ar.Timeout)
	defer cancel()

	query := `
		SELECT
			user_id,
			name,
			created_at
		FROM
			accounts
		ORDER BY
			user_id
	`

	rows, err := ar.DbSlave.QueryContext(ctx, query)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	accounts := make([]*model.Account, 0)
	for rows.Next() {
		var (
			aUserID    sql.NullInt64
			aName      sql.NullString
			aCreatedAt pq.NullTime
		)

		if err := rows.Scan(
			&aUserID,
			&aName,
			&aCreatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		account := model.Account{
			UserID:    aUserID.Int64,
			Name:      aName.String,
			CreatedAt: aCreatedAt.Time,
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (ar *accountRepository) Find(userID int64) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ar.Timeout)
	defer cancel()

	query := `
		SELECT
			user_id,
			name,
			created_at
		FROM
			accounts
		WHERE
			user_id = $1
	`

	var (
		aUserID    sql.NullInt64
		aName      sql.NullString
		aCreatedAt pq.NullTime
	)

	err := ar.DbSlave.QueryRowContext(ctx, query, userID).Scan(
		&aUserID,
		&aName,
		&aCreatedAt,
	)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.AccountNotExists
	}

	if err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	account := model.Account{
		UserID:    aUserID.Int64,
		Name:      aName.String,
		CreatedAt: aCreatedAt.Time,
	}

	return &account, nil
}

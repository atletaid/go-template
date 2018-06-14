package account

import (
	"github.com/sportivaid/go-template/src/model"
)

type AccountRepository interface {
	FindAll() ([]*model.Account, error)
	Find(userID int64) (*model.Account, error)
}

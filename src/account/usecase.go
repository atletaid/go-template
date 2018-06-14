package account

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sportivaid/go-template/src/model"
)

type Usecase interface {
	CreateToken(authUser *model.AuthUser, expTime time.Duration) (string, error)
	GetAccounts() ([]*model.Account, error)
	GetAccount(userID int64) (*model.Account, error)
}

type usecase struct {
	accountRepo AccountRepository
}

func NewAccountUsecase(
	accountRepo AccountRepository,
) Usecase {
	return &usecase{
		accountRepo: accountRepo,
	}
}

func (u *usecase) CreateToken(authUser *model.AuthUser, expTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  authUser.UserID,
		"username": authUser.Username,
		"email":    authUser.Email,
		"fullname": authUser.Fullname,
		"exp":      time.Now().Add(time.Second * expTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func (u *usecase) GetAccounts() ([]*model.Account, error) {
	users, err := u.accountRepo.FindAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (u *usecase) GetAccount(userID int64) (*model.Account, error) {
	user, err := u.accountRepo.Find(userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

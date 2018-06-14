package auth

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
)

type AuthUsecase interface {
	AuthenticateUserToken(token string) (*model.AuthUser, error)
}

type usecase struct{}

func (u *usecase) AuthenticateUserToken(tokenStr string) (*model.AuthUser, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		log.Println(err)
		return nil, apperror.TokenIsExpired
	}

	if !token.Valid {
		log.Println(apperror.InvalidAuthToken)
		return nil, apperror.InvalidAuthToken
	}

	authUser := model.AuthUser{
		UserID:   int64(claims["user_id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		Fullname: claims["fullname"].(string),
	}

	return &authUser, nil
}

func NewUsecase() AuthUsecase {
	return &usecase{}
}

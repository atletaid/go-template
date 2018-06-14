package rest

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
	"github.com/sportivaid/go-template/util/httputil"
)

type loginRequest struct {
	LoginType  string `json:"login_type" form:"login_type"`
	FirebaseID string `json:"firebase_id" form:"firebase_id"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := loginRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		if req.LoginType == "firebase" {
			authUser := model.AuthUser{
				UserID:   10000,
				Username: "example_username",
				Email:    "example_email",
				Fullname: "example_fullname",
			}

			token, err := au.CreateToken(&authUser, model.DefaultExpirationUser)
			if err != nil {
				log.Println(err)
				processTime := time.Now().Sub(startTime).Seconds()
				httputil.WriteErrorResponse(c, processTime, err)
				return
			}

			resp := loginResponse{token}

			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteResponse(c, []string{"Success login"}, processTime, resp)
		} else if req.LoginType == "password" {

		} else {
			log.Println(apperror.LoginTypeNotExists)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, apperror.LoginTypeNotExists)
			return
		}
	}
}

func GetUsers(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accounts, err := au.GetAccounts()
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get accounts"}, processTime, accounts)
	}
}

func GetUser(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		account, err := au.GetAccount(userID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get accounts"}, processTime, account)
	}
}

package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/sportivaid/go-template/src/common/apperror"
)

type response struct {
	StatusCode  int         `json:"status_code"`
	Messages    []string    `json:"messages"`
	ProcessTime float64     `json:"process_time"`
	Data        interface{} `json:"data"`
}

func WriteResponse(c *gin.Context, messages []string, processTime float64, data interface{}) {
	c.JSON(
		http.StatusOK,
		response{
			StatusCode:  http.StatusOK,
			Messages:    messages,
			ProcessTime: processTime,
			Data:        data,
		},
	)
}

func WriteErrorResponse(c *gin.Context, processTime float64, err error) {
	errCode := apperror.GetErrorCodes[err]
	if errCode == apperror.NotFoundErrorCode {
		errCode = apperror.DefaultErrorCode
	}

	c.JSON(
		errCode.HTTPcode,
		response{
			StatusCode:  errCode.StatusCode,
			Messages:    []string{err.Error()},
			ProcessTime: processTime,
			Data:        nil,
		},
	)
}

func DecodeFormRequest(r *http.Request, req interface{}) error {
	decoder := form.NewDecoder()
	r.ParseForm()
	return decoder.Decode(&req, r.Form)
}

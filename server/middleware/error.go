package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errHTTPProvider interface {
	HTTPStatusCode() int
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"-"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := ErrorResponse{Message: http.StatusText(http.StatusInternalServerError), Code: http.StatusInternalServerError}
		c.Next()

		for _, err := range c.Errors {
			if errProvider, ok := err.Err.(errHTTPProvider); ok {
				response.Message = err.Error()
				c.JSON(errProvider.HTTPStatusCode(), response)
				return
			} else {
				if validationErrors, ok := err.Err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
					response.Message = http.StatusText(http.StatusUnprocessableEntity)
					response.Error = validationErrors.Error()
					response.Code = http.StatusUnprocessableEntity
				}

				c.JSON(response.Code, response)
				return
			}
		}
	}
}

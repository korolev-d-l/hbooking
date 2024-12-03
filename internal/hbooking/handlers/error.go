package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type serviceError struct {
	err error
	msg string
}

func (e serviceError) Error() string {
	return e.msg
}

func (e serviceError) Unwrap() error {
	return e.err
}

func ValidationError(err error) serviceError {
	return serviceError{msg: "validation failed: " + err.Error(), err: err}
}

func ValidationErrorStr(err string) serviceError {
	return ValidationError(fmt.Errorf(err))
}

func ValidationFailed(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(serviceError); !ok {
		return false
	}

	return Error(c, err, http.StatusUnprocessableEntity)
}

func BadRequest(c *gin.Context, err error) bool {
	return Error(c, err, http.StatusBadRequest)
}

func InternalServerError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	_ = c.Error(err)
	c.Status(http.StatusInternalServerError)
	return true
}

func Error(c *gin.Context, err error, code int) bool {
	if err == nil {
		return false
	}

	_ = c.Error(err)
	c.AbortWithStatusJSON(code, ErrorResponse{err.Error()})
	return true
}

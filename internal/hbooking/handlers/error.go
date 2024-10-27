package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ValidationFailed(c *gin.Context, cond bool, description string) bool {
	if !cond {
		return false
	}

	err := fmt.Errorf("validation failed: %s", description)
	return Error(c, err, http.StatusUnprocessableEntity)
}

func BadRequest(c *gin.Context, err error, description string) bool {
	if err == nil {
		return false
	}

	err = fmt.Errorf("%s: %w", description, err)
	return Error(c, err, http.StatusBadRequest)
}

func Error(c *gin.Context, err error, code int) bool {
	if err == nil {
		return false
	}

	_ = c.Error(err)
	c.AbortWithStatusJSON(code, ErrorResponse{err.Error()})
	return true
}

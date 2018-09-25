package users

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"../../models"
	"../../utils"
)

// UserHandler struct
type UserHandler struct {
	db  *sqlx.DB
	lck sync.RWMutex
}

// NewUserHandler return new UserHandler object
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// errorMessage return error message as json string
func errorMessage(status int, msg string) string {
	msgFinal := &models.ResponseMessage{Status: status, Message: msg, Info: "/docs/api/errors"}
	result, _ := utils.NewResultTransformer(msgFinal).ToJSON()
	return result
}

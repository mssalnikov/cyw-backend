package events

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"../../models"
	"../../utils"
)

// EventHandler struct
type EventHandler struct {
	db  *sqlx.DB
	lck sync.RWMutex
}

// NewUserHandler return new UserHandler object
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// errorMessage return error message as json string
func errorMessage(status int, msg string) string {
	msgFinal := &models.ResponseMessage{Status: status, Message: msg, Info: "/docs/api/errors"}
	result, _ := utils.NewResultTransformer(msgFinal).ToJSON()
	return result
}

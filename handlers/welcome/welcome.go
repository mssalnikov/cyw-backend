package welcome

import (
	"../../models"
	"../../utils"
)

// WelcomeHandler struct
type WelcomeHandler struct {
}

// NewUserHandler return new UserHandler object
func NewWelcomeHandler() *WelcomeHandler{
	return &WelcomeHandler{}
}

// errorMessage return error message as json string
func errorMessage(status int, msg string) string {
	msgFinal := &models.ResponseMessage{Status: status, Message: msg, Info: "/docs/api/errors"}
	result, _ := utils.NewResultTransformer(msgFinal).ToJSON()
	return result
}

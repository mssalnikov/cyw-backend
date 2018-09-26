package events

import (
	"net/http"
)

func (uh *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	//// Todo get userId
	//userId := int64(2)
	//user, err := uh.getUser(userId)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
	//	return
	//}
	//
	//result, err := user.ToJSON()
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//fmt.Fprint(w, result)
}

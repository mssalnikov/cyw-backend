package events

import (
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	u "../../utils"
	"strconv"
)

func (eh *EventHandler) NewEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	token := r.Header.Get("auth_token")
	userId, err := u.RedisCon.Get(fmt.Sprintf("TOKEN:%s", token)).Int64()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	var event NewEvent
	err = json.Unmarshal(body, &event)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	eventId, err := eh.createEvent(event, userId)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}


	if eventId == 0 {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	//points, err := eh.createPoints()
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


func (eh *EventHandler) AllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := eh.allEvents()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	//
	//result := new(bytes.Buffer)
	//
	//w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//json.NewEncoder(result).Encode(events)
	//w.Write(result.Bytes())
	result, err := events.ToJSON()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, result)
}

func (eh *EventHandler) MyEvents(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("auth_token")
	userId, err := u.RedisCon.Get(fmt.Sprintf("TOKEN:%s", token)).Int64()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	events, err := eh.myEvents(userId)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	result, err := events.ToJSON()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, result)
}


func (eh *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	u, err := eh.getEvent(i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	result, err := u.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, result)
}

func (eh *EventHandler) JoinEvent(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("auth_token")
	userId, err := u.RedisCon.Get(fmt.Sprintf("TOKEN:%s", token)).Int64()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	id := r.URL.Query().Get("id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	u, err := eh.joinEvent(i, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	result, err := u.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, result)
}

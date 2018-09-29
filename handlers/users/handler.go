package users

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
)


func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	u, err := uh.getUser(i)
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


func (uh *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	var fbAuth FbAuth
	err = json.Unmarshal(body, &fbAuth)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	if fbAuth.FbAccessToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	token, err := uh.auth(fbAuth.FbAccessToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	result, err := token.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, result)
}

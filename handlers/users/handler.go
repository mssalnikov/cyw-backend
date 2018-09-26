package users

import (
	"fmt"
	"net/http"
	"../../utils"
	"io/ioutil"
	"encoding/json"
	"log"
)


//// GetUsers return all users
//func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
//	getDefaultHeader(w)
//
//	users, err := uh.getUsers().ToJSON()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	fmt.Fprint(w, users)
//}

// GetUser return user by id
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	//userId, err := utils.TryToGetUserIdByToken(r.Header["X-Session-Token"][0])
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
	//	return
	//}
	//
	//// get user by id
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


// GetUser return user by id
func (uh *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	//userId, err := utils.TryToGetUserIdByToken(r.Header["X-Session-Token"][0])
	userId, err := utils.CheckAuthToken(r.Header["X-Session-Token"][0])

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	// get user by id
	user, err := uh.getUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	result, err := user.ToJSON()
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(token)
}

//// issueSession issues a cookie session after successful Facebook login
//func (uh *UserHandler) IssueSession(redirectUrl string) http.Handler {
//	fn := func(w http.ResponseWriter, req *http.Request) {
//		ctx := req.Context()
//		facebookUser, err := facebook.UserFromContext(ctx)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		log.Println(facebookUser)
//		// 2. Implement a success handler to issue some form of session
//		session := sessionStore.New(sessionName)
//		session.Values[sessionUserKey] = facebookUser.ID
//		session.Save(w)
//		//url := redirectUrl + "/profile"
//		//http.Redirect(w, req, url, http.StatusFound)
//	}
//	return http.HandlerFunc(fn)
//}

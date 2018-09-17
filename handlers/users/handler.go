package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../../utils"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

func getDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

var Guid = xid.New()

const (
	sessionName    = "facebook-session"
	sessionUserKey = "facebookID"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(Guid.String()), nil)

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
	getDefaultHeader(w)
	userId, err := utils.TryToGetUserIdByToken(r.Header["X-Session-Token"][0])
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
	fmt.Fprint(w, result)
}

//
//// UpdateUser update user
//func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
//	getDefaultHeader(w)
//
//	if r.Body == nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))
//		return
//	}
//
//	decoder := json.NewDecoder(r.Body)
//	defer r.Body.Close()
//
//	var user User
//
//	err := decoder.Decode(&user)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	err = uh.updateUser(user)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

//// ProfileHandler
//func ProfileHandler(w http.ResponseWriter, r *http.Request) {
//	getDefaultHeader(w)
//
//	if r.Body == nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))
//		return
//	}
//
//	decoder := json.NewDecoder(r.Body)
//	defer r.Body.Close()
//
//	var user User
//
//	err := decoder.Decode(&user)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	err = uh.updateUser(user)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

//// DeleteUser delete user
//func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
//	getDefaultHeader(w)
//
//	vars := mux.Vars(r)
//
//	id, err := strconv.ParseInt(vars["id"], 10, 64)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	err = uh.deleteUserByID(id)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

// InsertUser insert new user into database
func (uh *UserHandler) InsertUser(w http.ResponseWriter, r *http.Request) {
	getDefaultHeader(w)

	//oauth2Config := &oauth2.Config{
	//	ClientID:     config.FacebookClientID,
	//	ClientSecret: config.FacebookClientSecret,
	//	RedirectURL:  "http://localhost:8080/user/callback",
	//	Endpoint:     facebookOAuth2.Endpoint,
	//	Scopes:       []string{"email"},
	//}
	//stateConfig := gologin.DebugOnlyCookieConfig

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var user User

	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	_, err = uh.insertUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMessage(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Login via Facebook
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	getDefaultHeader(w)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	err = uh.deleteUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// issueSession issues a cookie session after successful Facebook login
func (uh *UserHandler) IssueSession(redirectUrl string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(facebookUser)
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = facebookUser.ID
		session.Save(w)
		//url := redirectUrl + "/profile"
		//http.Redirect(w, req, url, http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

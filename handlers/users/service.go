package users

import (
	"errors"

	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	u "../../utils"
	"database/sql"
	"github.com/satori/go.uuid"
	"github.com/bitly/go-simplejson"
	"github.com/gkiryaziev/go-gorilla-sqlx/models"
	"github.com/gkiryaziev/go-gorilla-sqlx/utils"
)


type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

type FBError struct {
	Message string `json:"message"`
	Type string `json:"type"`
	OAuthException string `json:"OAuthException"`
	Code int `json:"code"`
	Fbtrace_id string `json:"fbtrace_id"`
}

type AccessResponse struct {
	Name string `json:"name"`
	Id string  `json:"id"`
	Email string `json:"email"`
	Error FBError `json:"error,omitempty"`
}


// getUsers return all users from db
func (us *UserHandler) getUsers() *u.ResultTransformer {

	// concurrency safe
	us.lck.RLock()
	defer us.lck.RUnlock()

	users := []User{}

	err := us.db.Select(&users, "select * from tbl_users order by id")
	if err != nil {
		panic(err)
	}

	header := models.Header{Status: "ok", Count: len(users), Data: users}
	result := u.NewResultTransformer(header)

	return result
}

// getUser return user by id from db
func (us *UserHandler) getUser(id int64) (*utils.ResultTransformer, error) {
	// concurrency safe
	us.lck.RLock()
	defer us.lck.RUnlock()

	user := User{}

	err := us.db.Get(&user, "select * from users where id = ?", id)
	if err != nil {
		return nil, err
	}

	header := models.Header{Status: "ok", Count: 1, Data: user}
	result := utils.NewResultTransformer(header)

	return result, nil
}

// deleteUserByID delete user by id and get rows affected in db
func (us *UserHandler) deleteUserByID(id int64) error {
	// concurrency safe
	us.lck.Lock()
	defer us.lck.Unlock()

	result, err := us.db.NamedExec("delete from users where id = :id", map[string]interface{}{"id": id})
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows <= 0 {
		return errors.New("0 Rows Affected")
	}

	return nil
}

// deleteUser delete user and get rows affected in db
func (us *UserHandler) deleteUser(user User) error {

	// concurrency safe
	us.lck.Lock()
	defer us.lck.Unlock()

	result, err := us.db.NamedExec("delete from users where id = :id", user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows <= 0 {
		return errors.New("0 Rows Affected")
	}

	return nil
}

// insertUser insert new user and get last id from db
func (us *UserHandler) insertUser(user User) (int64, error) {
	// concurrency safe
	us.lck.Lock()
	defer us.lck.Unlock()

	result, err := us.db.NamedExec("insert into users (fb_uid, username, email) values (:fb_uid, :username, :email)", user)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (us *UserHandler) auth(fbAccessToken string) ([]byte, error) {
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", fbAccessToken)
	response, err := http.Get(url)

	if err != nil {

		fmt.Println(err)
		return nil, err
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var resp AccessResponse
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.Error.Message != "" {
		return nil, errors.New(resp.Error.Message)
	}

	var user User

	row := u.DBCon.QueryRow(`SELECT * FROM users WHERE fb_uid = $1`, resp.Id)
	err = row.Scan(&user.Id, &user.FBUid, &user.UserName, &user.Email)

	jResponse := simplejson.New()

	switch err {
	case sql.ErrNoRows:
		id := 0
		sqlStatement := `INSERT INTO users (fb_uid, username, email) VALUES ($1, $2, $3) RETURNING id`
		err = u.DBCon.QueryRow(sqlStatement, resp.Id, resp.Name, resp.Email).Scan(&id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println("New record ID is:", id)
		newToken := uuid.Must(uuid.NewV4())

		key := fmt.Sprintf("TOKEN:%s", newToken)
		u.RedisCon.Set(key, id, 0).Err()
		jResponse.Set("auth_token", newToken.String())
		jResponse.Set("username", resp.Name)

		payload, err := jResponse.MarshalJSON()
		if err != nil {
			log.Println(err)
		}
		return payload, nil
	case nil:
		//log.Println(user)
		newToken := uuid.Must(uuid.NewV4())
		//fmt.Printf("UUIDv4: %s\n", u1)
		key := fmt.Sprintf("TOKEN:%s", newToken)
		u.RedisCon.Set(key, user.Id, 0).Err()
		jResponse.Set("auth_token", newToken.String())
		jResponse.Set("username", user.UserName)

		payload, err := jResponse.MarshalJSON()
		if err != nil {
			log.Println(err)
		}
		return payload, nil
	default:
		log.Println("Smth went wrong")
		return nil, &errorString{"Smth went wrong"}
	}
}

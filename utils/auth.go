package utils

import (
	"github.com/go-redis/redis"
	"fmt"
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


func CheckAuthToken(token string) (int64, error) {
	key := fmt.Sprintf("TOKEN:%s", token)

	val, err := RedisCon.Get(key).Int64()
	if err == redis.Nil {
		return 0, &errorString{"Key not found"}
	} else if err != nil {
		fmt.Println(err)
		return 0, err
	} else {
		fmt.Println("key2", val)
		return val, nil
	}
}
//
//func Auth(fbAccessToken string) (string, error){
//
//	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", fbAccessToken)
//	response, err := http.Get(url)
//
//	if err != nil {
//		fmt.Printf("%s", err)
//		return "", err
//	}
//	contents, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Printf("%s", err)
//		return "", err
//	}
//
//	var resp AccessResponse
//	err = json.Unmarshal(contents, &resp)
//
//	if err != nil {
//		log.Println("err")
//		panic(err)
//	}
//	if resp.Error.Message != "" {
//		log.Println("err")
//	}
//
//	var user u.User
//
//	row := DBCon.QueryRow(`SELECT * FROM users WHERE fb_uid = $1 ORDER BY id`, resp.Id)
//	err = row.Scan(&user.Id, &user.FBUid, &user.UserName, &user.Email, &user.DateAdded)
//
//	switch err {
//	case sql.ErrNoRows:
//		id := 0
//		sqlStatement := `INSERT INTO users (fb_uid, username, email) VALUES ($1, $2, $3, $4) RETURNING id`
//		err = DBCon.QueryRow(sqlStatement, resp.Id, resp.Name, resp.Email).Scan(&id)
//		if err != nil {
//			log.Println("Smth went wrong")
//			return "", &errorString{"Smth went wrong"}
//		}
//		log.Println("New record ID is:", id)
//		newToken := uuid.Must(uuid.NewV4())
//		//fmt.Printf("UUIDv4: %s\n", u1)
//		key := fmt.Sprintf("TOKEN:%s", newToken)
//		redisClient.Set(key, id, 0).Err()
//		return newToken.String(), nil
//	case nil:
//		log.Println(user)
//		newToken := uuid.Must(uuid.NewV4())
//		//fmt.Printf("UUIDv4: %s\n", u1)
//		key := fmt.Sprintf("TOKEN:%s", newToken)
//		redisClient.Set(key, user.Id, 0).Err()
//		return newToken.String(), nil
//	default:
//		log.Println("Smth went wrong")
//		return "", &errorString{"Smth went wrong"}
//	}
//
//}
//

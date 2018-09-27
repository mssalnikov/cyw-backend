package events

import (
	"../../conf"
	"log"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"github.com/simplereach/timeutils"
	u "../../utils"
)

func (es *EventHandler) createEvent(event NewEvent, userId int64) (int64, error) {
	// concurrency safe
	es.lck.Lock()
	defer es.lck.Unlock()
type AcceptEvent struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NewEvent struct {
	Points      []EventPoint   `json:"points"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Start       timeutils.Time `json:"start"`
	Finish      timeutils.Time `json:"finish"`
	IsPrivate   bool           `json:"is_private"`
}

type EventPoint struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
}
	var id int64
	sqlStatement := `INSERT INTO events (user_id, name, description, start, finish, is_private) values ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := u.DBCon.QueryRow(sqlStatement, userId, event.Name, event.Description, event.Start.Time, event.Finish.Time, event.IsPrivate).Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}
	for _, point := range event.Points {
		go es.createPoint(point, id)
	}
	return id, nil
}

func (es *EventHandler) createPoint(point EventPoint, eventId int64) error {
	// concurrency safe
	es.lck.Lock()
	defer es.lck.Unlock()

	err := es.createNaviaddress(point)
	if err != nil {
		return err
	}
	//result, err := es.db.Exec("insert into events (user_id, name, description, start_at, finish_at, is_private) values (?, ?, ?, ?, ?, ?)",
	//	userId, event.Name, event.Description, event.StartAt, event.FinishAt, event.IsPrivate)
	//
	//if err != nil {
	//	return 0, err
	//}
	//
	//id, err := result.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}

	return nil
}

func (es *EventHandler) createNaviaddress(point EventPoint) error {
	config, err := conf.NewConfig("config.yaml").Load()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%saddresses", config.Navi.ApiUri)
	fmt.Println("URL:>", url)

	naviPost := NaviaddressPost{Lat:point.Lat, Lng:point.Lng, DefaultLang:"ru", AddressType:"free"}

	b, err := json.Marshal(naviPost)
	if err != nil {
		return err
	}
	log.Println(naviPost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("auth-token", config.Navi.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return nil
}

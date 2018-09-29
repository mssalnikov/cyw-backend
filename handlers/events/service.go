package events

import (
	"../../conf"
	"log"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	u "../../utils"
)

func (es *EventHandler) allEvents(event NewEvent) ([]Event, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT id, user_id, name, description, start, finish FROM events WHERE is_finished < NOW()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	got := []Event{}
	for rows.Next() {
		var r Event
		err = rows.Scan(&r.Id, &r.UserId, &r.Description, &r.Start, &r.Finish)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
		got = append(got, r)
	}
	return got, nil
}

func (es *EventHandler) myEvents(event NewEvent, userId int64) ([]Event, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT id, user_id, name, description, start, finish FROM events WHERE user_id = ?", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	got := []Event{}
	for rows.Next() {
		var r Event
		err = rows.Scan(&r.Id, &r.UserId, &r.Description, &r.Start, &r.Finish)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
		got = append(got, r)
	}
	return got, nil
}

func (es *EventHandler) joinEvent(eventId int64, userId int64) (bool, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var id int64
	sqlStatement := `INSERT INTO userevent (user_id, event_id) values ($1, $2) RETURNING id`
	err := u.DBCon.QueryRow(sqlStatement, userId, eventId).Scan(&id)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}


func (es *EventHandler) checkCode(userId int64, pointId int64, token int64) (bool, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var rightToken int64
	err := u.DBCon.QueryRow("SELECT token FROM points where point_id = ?", pointId).Scan(&rightToken)
	if err != nil {
		log.Fatal(err)
	}

	if rightToken == token {
		return true, nil
	}
	return false, nil
}

func (es *EventHandler) answerQuestion(userId int64, pointId int64, answer string) (bool, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var rightAnswer string
	err := u.DBCon.QueryRow("SELECT answer FROM points where point_id = ?", pointId).Scan(&rightAnswer)
	if err != nil {
		log.Fatal(err)
	}

	if rightAnswer == answer {
		var id int64
		sqlStatement := `INSERT INTO userpoint (user_id, event_id, is_solved) values ($1, $2, $3) RETURNING id`
		err := u.DBCon.QueryRow(sqlStatement, userId, pointId, true).Scan(&id)

		if err != nil {
			log.Println(err)
			return false, err
		}

		if err != nil {
			log.Println(err)
			return false, err
		}


		return true, nil
	}
	return false, nil
}

func (es *EventHandler) joinedEvents(event NewEvent, userId int64) ([]Event, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT ev.id, ev.user_id, ev.name, ev.description, ev.start, ev.finish FROM events as ev " +
		"LEFT OUTER JOIN userevent as uev on ev.id = uev.event_id WHERE uev.user_id = ?", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	got := []Event{}
	for rows.Next() {
		var r Event
		err = rows.Scan(&r.Id, &r.UserId, &r.Description, &r.Start, &r.Finish)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
		got = append(got, r)
	}
	return got, nil
}

func (es *EventHandler) createEvent(event NewEvent, userId int64) (int64, error) {
	// concurrency safe
	es.lck.Lock()
	defer es.lck.Unlock()

	var id int64
	sqlStatement := `INSERT INTO events (user_id, name, description, start, finish) values ($1, $2, $3, $4, $5) RETURNING id`
	err := u.DBCon.QueryRow(sqlStatement, userId, event.Name, event.Description, event.Start.Time, event.Finish.Time).Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}

	// I love this async <3
	go es.createPoint(event.Points, id, userId)
	return id, nil
}

func (es *EventHandler) createPoint(points []EventPoint, eventId int64, userId int64) error {
	// concurrency safe
	es.lck.Lock()
	defer es.lck.Unlock()
	var prevPointId int64
	for _, point := range points {
			id, err := es.createNaviaddress(point, eventId, prevPointId)
			if err != nil {
				return err
			}
			prevPointId = id
			log.Println(prevPointId)
	}
	return nil
}

func (es *EventHandler) createNaviaddress(point EventPoint, eventId int64, prevPointId int64) (int64, error) {
	config, err := conf.NewConfig("config.yaml").Load()
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf("%saddresses", config.Navi.ApiUri)
	fmt.Println("URL:>", url)

	naviPost := NaviaddressPost{Lat:point.Lat, Lng:point.Lng, DefaultLang:"en", AddressType:"free"}

	b, err := json.Marshal(naviPost)
	if err != nil {
		return 0, err
	}
	log.Println(naviPost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("auth-token", config.Navi.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	var naviResp NaviaddressPostResponse
	err = json.Unmarshal(body, &naviResp)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var id int64
	sqlStatement := `INSERT INTO points (event_id, container, naviaddress, question, answer, token, prev_point_id) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = u.DBCon.QueryRow(sqlStatement, eventId, naviResp.Result.Container, naviResp.Result.Naviaddress, point.Question, point.Answer,
		u.RandomFourDigits(), prevPointId).Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	url = fmt.Sprintf("%saddresses/accept/%s/%s", config.Navi.ApiUri, naviResp.Result.Container, naviResp.Result.Naviaddress)
	fmt.Println("URL:>", url)
	naviAccept := NaviaddressAccept{Container:naviResp.Result.Container, Naviaddress:naviResp.Result.Naviaddress}

	b, err = json.Marshal(naviAccept)
	if err != nil {
		return 0, err
	}
	log.Println(naviAccept)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("auth-token", config.Navi.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))



	url = fmt.Sprintf("%saddresses/%s/%s", config.Navi.ApiUri, naviResp.Result.Container, naviResp.Result.Naviaddress)
	fmt.Println("URL:>", url)
	naviPut := NaviaddressPut{Name:"Cyw point", Description:point.Question, DefaultLang:"en", Lang:"en"}

	b, err = json.Marshal(naviPut)
	if err != nil {
		return 0, err
	}
	log.Println(naviPut)
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("auth-token", config.Navi.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return id, nil
}

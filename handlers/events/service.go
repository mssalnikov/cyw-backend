package events

import (
	"../../conf"
	"log"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"

	"encoding/json"
	u "../../utils"
	"time"
	"github.com/gkiryaziev/go-gorilla-sqlx/models"
	"github.com/gkiryaziev/go-gorilla-sqlx/utils"
)

func (es *EventHandler) allEvents() (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT id, user_id, name, description, start, finish FROM events WHERE finish < NOW()")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	got := []EventFromDB{}
	for rows.Next() {
		var r EventFromDB
		err = rows.Scan(&r.Id, &r.UserId, &r.Name, &r.Description, &r.Start, &r.Finish)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
		got = append(got, r)
	}
	header := models.Header{Status: "ok", Count: len(got), Data: got}
	result := utils.NewResultTransformer(header)
	return result, nil
}

func (es *EventHandler) myEvents(userId int64) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT id, user_id, name, description, start, finish FROM events WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	got := []EventFromDB{}
	for rows.Next() {
		var r EventFromDB
		err = rows.Scan(&r.Id, &r.UserId, &r.Name, &r.Description, &r.Start, &r.Finish)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
		got = append(got, r)
	}

	header := models.Header{Status: "ok", Count: len(got), Data: got}
	result := utils.NewResultTransformer(header)
	return result, nil
}


func (es *EventHandler) getPoint(userId int64, pointId int64) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var (
		id uint64
		name string
		question string
		token string
		isFound bool
		isSolved bool
	)
	err := u.DBCon.QueryRow("SELECT p.id, p.name, p.question, p.token, up.is_found, up.is_solved FROM points p LEFT OUTER JOIN userpoint as up on p.id = up.point_id WHERE p.id = $1", pointId).Scan(&id, &name, &question, &token, &isFound, &isSolved)
	if err != nil {
		log.Println(err)
	}

	res := PointFromDbForUser {
		Id:id,
		Name:name,
		Question:question,
		IsFound:isFound,
		IsSolved:isSolved,
	}
	header := models.Header{Status: "ok", Count: 1, Data: res}
	result := utils.NewResultTransformer(header)

	return result, nil
}

func (es *EventHandler) getEvent(eventId int64) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var (
		id uint64
		userId uint64
		name string
		description string
		start time.Time
		finish time.Time
	)
	err := u.DBCon.QueryRow("SELECT id, user_id, name, description, start, finish FROM events WHERE id = $1", eventId).Scan(&id, &userId, &name, &description, &start, &finish)
	if err != nil {
		log.Println(err)
	}
	query := fmt.Sprintf("SELECT p.id, p.name, up.is_solved, up.is_found FROM points as p LEFT OUTER JOIN userpoint as up on p.id = up.point_id WHERE p.event_id = %d", eventId)
	rows, err := u.DBCon.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	points := []LabelPoints{}
	for rows.Next() {
		var p LabelPoints
		err = rows.Scan(&p.Id, &p.Name, &p.IsSolved, &p.IsFound)
		if err != nil {
			log.Printf("Scan: %v", err)
		}
		points = append(points, p)
	}
	res := EventFromDBForUser {
		Id:id,
		UserId:userId,
		Name:name,
		Description:description,
		Start:start,
		Finish:finish,
		Points:points,
	}
	header := models.Header{Status: "ok", Count: 1, Data: res}
	result := utils.NewResultTransformer(header)

	return result, nil
}

func (es *EventHandler) joinEvent(eventId int64, userId int64) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var id int64
	sqlStatement := `INSERT INTO userevent (user_id, event_id) values ($1, $2) RETURNING id`
	err := u.DBCon.QueryRow(sqlStatement, userId, eventId).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := fmt.Sprintf("SELECT id FROM points WHERE event_id = %d", eventId)
	rows, err := u.DBCon.Query(query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("can't find event")
	}
	defer rows.Close()

	for rows.Next() {
		var p PointId
		err = rows.Scan(&p.Id)
		if err != nil {
			log.Printf("Scan: %v", err)
			return nil, err
		}
			var upId int64
			sqlStatement := `INSERT INTO userpoint (user_id, point_id, is_solved, is_found) VALUES ($1, $2, $3, $4) RETURNING id`
			err := u.DBCon.QueryRow(sqlStatement, userId, p.Id, false, false).Scan(&upId)
			log.Println(upId)
			if err != nil {
				log.Println(err)
				return nil, errors.New("can't create userpoint")
			}
	}

	header := models.Header{Status: "ok"}
	result := utils.NewResultTransformer(header)

	return result, nil
}


func (es *EventHandler) checkCode(userId int64, pointId int64, token string) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var rightToken string
	err := u.DBCon.QueryRow("SELECT token FROM points where point_id = $1", pointId).Scan(&rightToken)
	if err != nil {
		log.Println(err)
		return nil, errors.New("can't find token")
	}
	if token == rightToken {
		sqlStatement := `UPDATE userpoint  SET is_found = TRUE WHERE user_id = $1 and point_id = $2;`
		_, err = u.DBCon.Exec(sqlStatement, userId, pointId)
		if err != nil {
			log.Println(err)
			return nil, errors.New("can't update userpoint")
		}

		header := models.Header{Status: "ok"}
		result := utils.NewResultTransformer(header)
		return result, nil
	} else {
		if err != nil {
			log.Println(err)
			return nil, err
		}

		header := models.Header{Status: "fail"}
		result := utils.NewResultTransformer(header)
		return result, errors.New("wrong token")
	}
}

func (es *EventHandler) checkEventSolved(userId int64, eventId int64) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var allPoints int64
	err := u.DBCon.QueryRow("SELECT COUNT(id) FROM points  WHERE event_id = $1", eventId).Scan(&allPoints)


	if err != nil {
		log.Println(err)
		return nil, err
	}

	var solvedPoints int64
	err = u.DBCon.QueryRow("SELECT COUNT(p.id) FROM points p left outer join userpoint up on p.id = up.point_id WHERE p.event_id = $1 and is_solved =TRUE and up.user_id = $2;", eventId, userId).Scan(&solvedPoints)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	if allPoints == solvedPoints {
		sqlStatement := `UPDATE userevent SET is_passed = TRUE WHERE user_id = $1 and event_id = $2;`
		_, err = u.DBCon.Exec(sqlStatement, userId, eventId)
		if err != nil {
			log.Println(err)
			return nil, errors.New("can't update userevent")
		}

		header := models.Header{Status: "ok"}
		result := utils.NewResultTransformer(header)
		return result, nil
	} else {
		if err != nil {
			log.Println(err)
			return nil, err
		}

		header := models.Header{Status: "fail"}
		result := utils.NewResultTransformer(header)
		return result, errors.New("not passed")
	}
}

func (es *EventHandler) answerQuestion(userId int64, pointId int64, answer string) (*utils.ResultTransformer, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	var rightAnswer string
	err := u.DBCon.QueryRow("SELECT answer FROM points where id = $1", pointId).Scan(&rightAnswer)
	if err != nil {
		log.Println(err)
	}

	if rightAnswer == answer {
		//var id int64
		sqlStatement := `UPDATE userpoint SET is_solved = true WHERE user_id = $1 and point_id = $2`
		_, err := u.DBCon.Exec(sqlStatement, userId, pointId)

		if err != nil {
			log.Println(err)
			header := models.Header{Status: "fail"}
			result := utils.NewResultTransformer(header)
			return result, errors.New("can't update userpoint")
		}

		header := models.Header{Status: "ok"}
		result := utils.NewResultTransformer(header)
		return result, nil
	} else {
		header := models.Header{Status: "fail"}
		result := utils.NewResultTransformer(header)
		return result, errors.New("wrong answer")
	}
}

func (es *EventHandler) joinedEvents(event NewEvent, userId int64) ([]Event, error) {
	es.lck.Lock()
	defer es.lck.Unlock()

	rows, err := u.DBCon.Query("SELECT ev.id, ev.user_id, ev.name, ev.description, ev.start, ev.finish FROM events as ev " +
		"LEFT OUTER JOIN userevent as uev on ev.id = uev.event_id WHERE uev.user_id = $1", userId)
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

	// I love this async <3
	go es.createPoints(event.Points, id, userId)
	return id, nil
}

func (es *EventHandler) createPoints(points []EventPoint, eventId int64, userId int64) error {
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
	sqlStatement := `INSERT INTO points (event_id, container, naviaddress, name, question, answer, token, prev_point_id) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = u.DBCon.QueryRow(sqlStatement, eventId, naviResp.Result.Container, naviResp.Result.Naviaddress, point.Name, point.Question, point.Answer,
		point.Token, prevPointId).Scan(&id)

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
	naviPut := NaviaddressPut{Name:point.Name, Description:point.Question, DefaultLang:"en", Lang:"en"}

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

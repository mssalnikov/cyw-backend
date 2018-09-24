package events

import (
	"../../conf"
	"bytes"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

type NewEvent struct {
	Name string `json:"name"`
	Description string  `json:"description"`
}


func makeEvent(name string, description string) error{
	config, err := conf.NewConfig("config.yaml").Load()
	if err != nil {
		return err
	}

	ev	 := &NewEvent{Name: name, Description: description}
	e, err := json.Marshal(ev)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Println(name)
	resp, err := http.Post(fmt.Sprintf("%s/addresses", config.Navi.ApiUri), "application/json", bytes.NewBuffer(e))
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Println(resp)
	return nil
}

func main() {
	makeEvent("ololo", "mamka tvoya")
}
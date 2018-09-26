package events

import (
	"../../conf"
	"bytes"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

type AcceptEvent struct {
	Name string `json:"name"`
	Description string  `json:"description"`
}

type NewEvent struct {
	Lat float32 `json:"lat"`
	Lng float32  `json:"lng"`
	DefaultLang string  `json:"default_lang"`
	AddressType string  `json:"address_type"`
}


func makeEvent(lat float32, lng float32, name string, description string) error{
	config, err := conf.NewConfig("config.yaml").Load()
	if err != nil {
		return err
	}

	// for our testing application we're making only free addresses in ru_RU locale
	ev := &NewEvent{
		Lat: lat,
		Lng: lng,
		DefaultLang: "ru",
		AddressType: "free"}

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

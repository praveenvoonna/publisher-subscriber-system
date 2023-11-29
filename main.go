package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"gopkg.in/antage/eventsource.v1"
)

type Object struct {
	TimeStamp int64 `json:"timestamp"`
	Data      Data  `json:"data"`
}

type Data struct {
	UserName string `json:"username"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
}

type ResponseData struct {
	Data struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	} `json:"data"`
	Support struct {
		URL  string `json:"url"`
		Text string `json:"text"`
	} `json:"support"`
}

var ResponseDataObj ResponseData

func GetDataFromSource() error {
	resp, err := http.Get("https://reqres.in/api/users/2")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &ResponseDataObj)
	if err != nil {
		return err
	}

	return nil
}

func startScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)
	interval, err := strconv.Atoi(os.Getenv("GetDataFromServerInSec"))
	if err != nil {
		log.Fatalf("Failed to parse GetDataFromServerInSec: %s\n", err)
	}
	_, err = s.Every(interval).Seconds().Do(func() {
		if err := GetDataFromSource(); err != nil {
			log.Printf("Error fetching data: %s\n", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to start scheduler: %s\n", err)
	}

	s.StartAsync()
	return s
}

func createObject() Object {
	emailPart := strings.Split(ResponseDataObj.Data.Email, "@")
	return Object{
		TimeStamp: time.Now().Unix(),
		Data: Data{
			UserName: ResponseDataObj.Data.FirstName + "." + ResponseDataObj.Data.LastName,
			Domain:   emailPart[1],
			Name:     ResponseDataObj.Data.FirstName + " " + ResponseDataObj.Data.LastName,
		},
	}
}

func sendEvents(es eventsource.EventSource) {
	id := 1
	for {
		obj := createObject()
		es.SendEventMessage(fmt.Sprintf("%+v", obj), "tick-event", strconv.Itoa(id))
		id++
		interval, err := strconv.Atoi(os.Getenv("PushDataToClientInSec"))
		if err != nil {
			log.Printf("Failed to parse PushDataToClientInSec: %s\n", err)
			interval = 5 // Default interval if parsing fails
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	err = GetDataFromSource()
	if err != nil {
		log.Fatalf("Error getting initial data: %s\n", err)
	}

	es := eventsource.New(nil, nil)
	defer es.Close()

	go sendEvents(es)
	startScheduler()

	http.Handle("/events", es)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

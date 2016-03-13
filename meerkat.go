package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Broadcaster struct {
	Id          string
	Name        string
	DisplayName string
}

type Broadcast struct {
	Id            string
	Broadcaster   Broadcaster
	Caption       string
	Location      string
	WatchersCount int
	CommentsCount int
}

type CompactBroadcast struct {
	Url           string
	Caption       string
	Location      string
	WatchersCount int
	CommentsCount int
}

type Result struct {
	Result []Broadcast
}

func getBroadcasts() ([]Broadcast, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://resources.meerkatapp.co/broadcasts", nil)
	if err != nil {
		log.Fatal("Could not create request object.")
	}
	req.Header.Add("Authorization", meerkatKey)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not perform request", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return result.Result, nil
}

func handleMeerkatBroadcasts(w http.ResponseWriter, r *http.Request) {
	broadcasts, err := getBroadcasts()
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return
	}
	var cbs []CompactBroadcast

	for _, b := range broadcasts {
		url := fmt.Sprintf("http://meerkatapp.co/%s/%s", b.Broadcaster.Name, b.Id)
		cb := CompactBroadcast{Url: url, Caption: b.Caption, Location: b.Location, WatchersCount: b.WatchersCount, CommentsCount: b.CommentsCount}
		cbs = append(cbs, cb)
	}
	b, err := json.Marshal(cbs)
	if err != nil {
		fmt.Fprintf(w, "JSON Error %v", err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

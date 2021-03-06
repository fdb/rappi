package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Url           string `json:"url"`
	Caption       string `json:"caption"`
	Location      string `json:"location"`
	WatchersCount int    `json:"watchersCount"`
	CommentsCount int    `json:"commentsCount"`
}

type MeerkatOk struct {
	Status     string             `json:"status"`
	Broadcasts []CompactBroadcast `json:"broadcasts"`
}

type Result struct {
	Result []Broadcast
}

func getBroadcasts() ([]Broadcast, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://resources.meerkatapp.co/broadcasts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", meerkatKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func handleMeerkatBroadcasts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	broadcasts, err := getBroadcasts()
	if err != nil {
		jsonError(w, err)
		return
	}
	var cbs []CompactBroadcast

	for _, b := range broadcasts {
		url := fmt.Sprintf("http://meerkatapp.co/%s/%s", b.Broadcaster.Name, b.Id)
		cb := CompactBroadcast{Url: url, Caption: b.Caption, Location: b.Location, WatchersCount: b.WatchersCount, CommentsCount: b.CommentsCount}
		cbs = append(cbs, cb)
	}
	b, err := json.Marshal(MeerkatOk{Status: "ok", Broadcasts: cbs})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BroadcastIn struct {
	Id              string
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	UserId          string `json:"user_id"`
	UserDisplayName string `json:"user_display_name"`
	Username        string
	TwitterUsername string `json:"twitter_username"`
	ProfileImageUrl string `json:"profile_image_url"`
	State           string
	Language        string
	Start           string
	City            string
	Country         string
	ImageUrl        string `json:"image_url"`
	Status          string
}

type BroadcastOut struct {
	Id              string `json:"id"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	UserId          string `json:"userId"`
	UserDisplayName string `json:"userDisplayName"`
	Username        string `json:"username"`
	TwitterUsername string `json:"twitterUsername"`
	ProfileImageUrl string `json:"profileImageUrl"`
	State           string `json:"state"`
	Language        string `json:"language"`
	Start           string `json:"start"`
	City            string `json:"city"`
	Country         string `json:"country"`
	ImageUrl        string `json:"imageUrl"`
	Status          string `json:"status"`
}

type BroadcastDetailsIn struct {
	Id        string
	Type      string
	HlsUrl    string `json:"hls_url"`
	ReplayUrl string `json:"replay_url"`
}

type BroadcastDetailsOut struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	StreamingUrl string `json:"streamingUrl"`
	ReplayUrl    string `json:"replayUrl"`
}

type PeriscopeOk struct {
	Status     string         `json:"status"`
	Broadcasts []BroadcastOut `json:"broadcasts"`
}

func periscopeGetBroadcasts() ([]BroadcastIn, error) {
	var jsonStr = []byte(`{"cookie":"` + periscopeCookie + `"}`)
	req, err := http.NewRequest("POST", "https://api.periscope.tv/api/v2/rankedBroadcastFeed", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var broadcasts []BroadcastIn
	err = json.Unmarshal(body, &broadcasts)
	if err != nil {
		return nil, err
	}
	return broadcasts, nil
}

func periscopeGetBroadcastDetails(broadcastId string) (BroadcastDetailsIn, error) {
	res, err := http.Get("https://api.periscope.tv/api/v2/getAccessPublic?broadcast_id=" + broadcastId)
	if err != nil {
		return BroadcastDetailsIn{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var details BroadcastDetailsIn
	err = json.Unmarshal(body, &details)
	if err != nil {
		return BroadcastDetailsIn{}, err
	}
	return details, nil
}

func handlePeriscopeBroadcasts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	broadcasts, err := periscopeGetBroadcasts()
	if err != nil {
		jsonError(w, err)
		return
	}
	var broadcastsOut []BroadcastOut

	for _, broadcast := range broadcasts {
		b := BroadcastOut{
			Id:              broadcast.Id,
			CreatedAt:       broadcast.CreatedAt,
			UpdatedAt:       broadcast.UpdatedAt,
			UserId:          broadcast.UserId,
			UserDisplayName: broadcast.UserDisplayName,
			Username:        broadcast.Username,
			TwitterUsername: broadcast.TwitterUsername,
			ProfileImageUrl: broadcast.ProfileImageUrl,
			State:           broadcast.State,
			Language:        broadcast.Language,
			Start:           broadcast.Start,
			City:            broadcast.City,
			Country:         broadcast.Country,
			ImageUrl:        broadcast.ImageUrl,
			Status:          broadcast.Status,
		}
		broadcastsOut = append(broadcastsOut, b)
	}
	b, err := json.Marshal(PeriscopeOk{Status: "ok", Broadcasts: broadcastsOut})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func handlePeriscopeBroadcastDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.FormValue("id")
	if id == "" {
		jsonError(w, errors.New("A broadcast id is required."))
		return
	}
	details, err := periscopeGetBroadcastDetails(id)
	if err != nil {
		jsonError(w, err)
		return
	}
	var broadcastType string
	if details.Type == "StreamTypeWeb" {
		broadcastType = "live"
	} else if details.Type == "StreamTypeReplay" {
		broadcastType = "replay"
	} else {
		jsonError(w, errors.New("Broadcast not found."))
		return
	}
	d := BroadcastDetailsOut{
		Id:           id,
		Type:         broadcastType,
		StreamingUrl: details.HlsUrl,
		ReplayUrl:    details.ReplayUrl,
	}
	enc := json.NewEncoder(w)
	m := make(map[string]interface{})
	m["status"] = "ok"
	m["details"] = d
	err = enc.Encode(m)
	if err != nil {
		jsonError(w, err)
		return
	}
}

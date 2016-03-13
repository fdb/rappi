package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterStatus struct {
	Url           string   `json:"url"`
	Id            int64    `json:"id"`
	CreatedAt     string   `json:"createdAt"`
	Text          string   `json:"text"`
	UserId        int64    `json:"userId"`
	UserHandle    string   `json:"user"`
	FavoriteCount int      `json:"favoriteCount"`
	RetweetCount  int      `json:"retweetCount"`
	ImageUrls     []string `json:"imageUrls"`
}

type TwitterOk struct {
	Status string          `json:"status"`
	Tweets []TwitterStatus `json:"tweets"`
}

type TwitterError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var twitterApi *anaconda.TwitterApi

func initTwitter(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	twitterApi = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func twitterStatusUrl(status anaconda.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", status.User.ScreenName, status.IdStr)
}

func jsonError(w http.ResponseWriter, err error) {
	errorString := fmt.Sprintf("%v", err)
	fmt.Fprintf(w, `{"status":"error","message":%s}`, strconv.Quote(errorString))
}

func handleTwitterSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	searchResult, err := twitterApi.GetSearch(query, nil)
	if err != nil {
		jsonError(w, err)
		return
	}
	var statuses []TwitterStatus
	for _, status := range searchResult.Statuses {
		ts := TwitterStatus{
			Url:           twitterStatusUrl(status),
			Id:            status.Id,
			CreatedAt:     status.CreatedAt,
			Text:          status.Text,
			UserId:        status.User.Id,
			UserHandle:    status.User.ScreenName,
			FavoriteCount: status.FavoriteCount,
			RetweetCount:  status.RetweetCount,
		}
		ts.ImageUrls = make([]string, 0)
		for _, media := range status.Entities.Media {
			ts.ImageUrls = append(ts.ImageUrls, media.Media_url_https)
		}
		statuses = append(statuses, ts)
	}
	b, err := json.Marshal(TwitterOk{Status: "ok", Tweets: statuses})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type TwitterSearchResult struct {
	Query string
	Statuses []TwitterStatus
	Time time.Time
}

var twitterApi *anaconda.TwitterApi
var twitterCache map[string]*TwitterSearchResult

func initTwitter(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	twitterApi = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	twitterCache = make(map[string]*TwitterSearchResult)
}

func twitterStatusUrl(status anaconda.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", status.User.ScreenName, status.IdStr)
}

func doTwitterSearch(query string) (*TwitterSearchResult, error) {
	searchResult, err := twitterApi.GetSearch(query, nil)
	if err != nil {
		return nil, err
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

	return &TwitterSearchResult{query, statuses, time.Now()}, nil
}

func handleTwitterSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.FormValue("q")
	result1, ok := twitterCache[query]
	if (!ok || time.Now().Sub(result1.Time) > 5 * time.Minute) {
		fmt.Println("New search " + query)
		result2, err := doTwitterSearch(query)
		if err != nil {
			jsonError(w, err)
			return
		}
		twitterCache[query] = result2
		result1 = result2
	}

	// searchResult, err := twitterApi.GetSearch(query, nil)
	b, err := json.Marshal(TwitterOk{Status: "ok", Tweets: result1.Statuses})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

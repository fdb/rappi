package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	ImageUrls     []string `json:"imageUrls"`
}

type TwitterOk struct {
	Status string          `json:"status"`
	Tweets []TwitterStatus `json:"tweets"`
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

// https://developer.twitter.com/en/docs/tweets/search/api-reference/get-search-tweets.html
func handleTwitterSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.FormValue("q")
	v := url.Values{}
	v.Set("geocode", r.FormValue("geocode")) // 48.858278,2.294254,10km
	v.Set("count", r.FormValue("count")) // 15 by default, max 100
	v.Set("result_type", r.FormValue("result_type")) // mixed, recent, popular
	v.Set("lang", r.FormValue("lang")) // ISO 639-1 code: en,fr,es,... (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
	v.Set("until", r.FormValue("until")) // tweets before data, in YYYY-MM-DD format. 7-day limit

	searchResult, err := twitterApi.GetSearch(query, v)
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
		ts.Latitude, _ = status.Latitude()
		ts.Longitude, _ = status.Longitude()
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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	PlaceID       string   `json:"placeId"`
	Place         string   `json:"place"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
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

// https://developer.twitter.com/en/docs/tweets/search/api-reference/get-search-tweets.html
func doTwitterSearch(r *http.Request) (*TwitterSearchResult, error) {
	query := r.FormValue("q")
	v := url.Values{}
	v.Set("geocode", r.FormValue("geocode")) // 48.858278,2.294254,10km
	v.Set("count", r.FormValue("count")) // 15 by default, max 100
	v.Set("result_type", r.FormValue("result_type")) // mixed, recent, popular
	v.Set("lang", r.FormValue("lang")) // ISO 639-1 code: en,fr,es,... (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
	v.Set("until", r.FormValue("until")) // tweets before data, in YYYY-MM-DD format. 7-day limit
	searchResult, err := twitterApi.GetSearch(query, v)
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
			PlaceID:       status.Place.ID,
			Place:         status.Place.Name,
		}
		if status.HasCoordinates() {
			ts.Latitude, _ = status.Latitude()
			ts.Longitude, _ = status.Longitude()
		} else if len(status.Place.ID) > 0 {
			bbox := status.Place.BoundingBox.Coordinates[0]
			if len(bbox) == 4 {
				lon := (bbox[0][0] + bbox[1][0] + bbox[2][0] + bbox[3][0]) / 4.0
				lat := (bbox[0][1] + bbox[1][1] + bbox[2][1] + bbox[3][1]) / 4.0
				ts.Latitude = lon
				ts.Longitude = lat
			}
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

	query := r.URL.RawQuery
	result1, ok := twitterCache[query]
	if (!ok || time.Now().Sub(result1.Time) > 5 * time.Minute) {
		fmt.Println("New search " + query)
		result2, err := doTwitterSearch(r)
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

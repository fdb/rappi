package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
)

var twitterApi *anaconda.TwitterApi

func initTwitter(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	twitterApi = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func handleTwitterSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	searchResult, err := twitterApi.GetSearch(query, nil)
	if err != nil {
		fmt.Fprintf(w, "Twitter API error: %v", err)
		return
	}
	b, err := json.Marshal(searchResult.Statuses)
	if err != nil {
		fmt.Fprintf(w, "Twitter JSON Error: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

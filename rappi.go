package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"text/template"
)

func jsonError(w http.ResponseWriter, err error) {
	errorString := fmt.Sprintf("%v", err)
	fmt.Fprintf(w, `{"status":"error","message":%s}`, strconv.Quote(errorString))
}

func renderTemplate(w http.ResponseWriter, name string) {
	lp := path.Join("templates", "_base.html")
	fp := path.Join("templates", name)
	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "base", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
}

func handleTwitterIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "twitter.html")
}

func handleMeerkatIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "meerkat.html")
}

var meerkatKey string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)

	// Twitter
	twitterConsumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if twitterConsumerKey == "" || twitterConsumerSecret == "" || twitterAccessToken == "" || twitterAccessTokenSecret == "" {
		log.Fatal("Twitter: check if TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_SECRET, TWITTER_ACCESS_TOKEN and TWITTER_ACCESS_TOKEN_SECRET environement variables are set. Go to https://apps.twitter.com/ to obtain keys.")
	}

	initTwitter(twitterConsumerKey, twitterConsumerSecret, twitterAccessToken, twitterAccessTokenSecret)
	http.HandleFunc("/twitter/", handleTwitterIndex)
	http.HandleFunc("/twitter/search.json", handleTwitterSearch)

	// Meerkat
	meerkatKey = os.Getenv("MEERKAT_KEY")
	if meerkatKey == "" {
		log.Fatal("No $MEERKAT_KEY set.")
	}
	http.HandleFunc("/meerkat/", handleMeerkatIndex)
	http.HandleFunc("/meerkat/broadcasts.json", handleMeerkatBroadcasts)

	fmt.Println("http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

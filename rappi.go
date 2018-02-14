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

func handleFlickrIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "flickr.html")
}

func handleImgurIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "imgur.html")
}

func handlePeriscopeIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "periscope.html")
}

func handleMeerkatIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "meerkat.html")
}

var flickrKey string
var imgurClientId string
var periscopeCookie string
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
	if len(twitterConsumerKey) > 0 && len(twitterConsumerSecret) > 0 && len(twitterAccessToken) > 0 && len(twitterAccessTokenSecret) > 0 {
		initTwitter(twitterConsumerKey, twitterConsumerSecret, twitterAccessToken, twitterAccessTokenSecret)
		http.HandleFunc("/twitter/", handleTwitterIndex)
		http.HandleFunc("/twitter/search.json", handleTwitterSearch)
	} else {
		log.Println("Twitter: check if TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_SECRET, TWITTER_ACCESS_TOKEN and TWITTER_ACCESS_TOKEN_SECRET environement variables are set. Go to https://apps.twitter.com/ to obtain keys.")
	}

	// Flickr
	flickrKey = os.Getenv("FLICKR_KEY")
	if len(flickrKey) > 0 {
		http.HandleFunc("/flickr/", handleFlickrIndex)
		http.HandleFunc("/flickr/search.json", handleFlickrSearch)
	} else {
		log.Println("Flickr: no FLICKR_KEY environment variable found. Obtain API key here: https://www.flickr.com/services/api/keys/apply/")
	}

	// Imgur
	imgurClientId = os.Getenv("IMGUR_CLIENT_ID")
	if len(imgurClientId) > 0 {
		http.HandleFunc("/imgur/", handleImgurIndex)
		http.HandleFunc("/imgur/search.json", handleImgurSearch)
	} else {
		log.Println("Imgur: no IMGUR_CLIENT_ID environment variable found. Obtain API key here: https://api.imgur.com/oauth2/addclient")
	}

	// Periscope
	periscopeCookie = os.Getenv("PERISCOPE_COOKIE")
	if len(periscopeCookie) > 0  {
		http.HandleFunc("/periscope/", handlePeriscopeIndex)
		http.HandleFunc("/periscope/broadcasts.json", handlePeriscopeBroadcasts)
		http.HandleFunc("/periscope/broadcastDetails.json", handlePeriscopeBroadcastDetails)
	} else {
		log.Println("Periscope: no PERISCOPE_COOKIE environment variable found.")
	}

	// Meerkat
	meerkatKey = os.Getenv("MEERKAT_KEY")
	if len(meerkatKey) > 0 {
		http.HandleFunc("/meerkat/", handleMeerkatIndex)
		http.HandleFunc("/meerkat/broadcasts.json", handleMeerkatBroadcasts)
	} else {
		log.Println("No $MEERKAT_KEY set.")
	}

	fmt.Println("http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, name string) {
	lp := path.Join("templates", "_base.html")
	fp := path.Join("templates", name)
	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "base", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ImgurPhotoIn struct {
	Id          string
	Title       string
	Description string
	Datetime    int64
	Link        string
	IsAlbum     bool   `json:"is_album"`
	AccountUrl  string `json:"account_url"`
}

type ImgurResult struct {
	Data []ImgurPhotoIn
}

type ImgurPhotoOut struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    int64  `json:"datetime"`
	ImageUrl    string `json:"imageUrl"`
	UserHandle  string `json:"userHandle"`
}

type ImgurOk struct {
	Status string          `json:"status"`
	Photos []ImgurPhotoOut `json:"photos"`
}

func imgurSearch(query string) ([]ImgurPhotoIn, error) {
	client := &http.Client{}
	v := url.Values{}
	v.Set("q", query)
	req, err := http.NewRequest("GET", "https://api.imgur.com/3/gallery/search/time?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Client-ID "+imgurClientId)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var result ImgurResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func handleImgurSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.FormValue("q")
	photos, err := imgurSearch(query)
	if err != nil {
		jsonError(w, err)
		return
	}
	var photosOut []ImgurPhotoOut
	for _, photo := range photos {
		if !photo.IsAlbum {
			p := ImgurPhotoOut{
				Id:          photo.Id,
				Title:       photo.Title,
				Description: photo.Description,
				Datetime:    photo.Datetime,
				ImageUrl:    photo.Link,
				UserHandle:  photo.AccountUrl,
			}
			photosOut = append(photosOut, p)
		}
	}
	b, err := json.Marshal(ImgurOk{Status: "ok", Photos: photosOut})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

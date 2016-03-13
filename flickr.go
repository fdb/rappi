package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type FlickrPhotoIn struct {
	Id     string
	Owner  string
	Secret string
	Server string
	Farm   int
	Title  string
}

type FlickrPhotoListIn struct {
	Photo []FlickrPhotoIn
}

type FlickrPhotosIn struct {
	Photos FlickrPhotoListIn
}

type FlickrPhotoOut struct {
	Id           string `json:"id"`
	Owner        string `json:"owner"`
	Title        string `json:"title"`
	SquareUrl    string `json:"squareUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	MediumUrl    string `json:"mediumUrl"`
	LargeUrl     string `json:"largeUrl"`
}

type FlickrOk struct {
	Status string           `json:"status"`
	Photos []FlickrPhotoOut `json:"photos"`
}

func flickrSearch(query string) ([]FlickrPhotoIn, error) {
	v := url.Values{}
	v.Set("method", "flickr.photos.search")
	v.Set("api_key", flickrKey)
	v.Set("tags", query)
	v.Set("per_page", "100")
	v.Set("format", "json")
	v.Set("nojsoncallback", "1")
	res, err := http.Get("https://api.flickr.com/services/rest/?" + v.Encode())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	var result FlickrPhotosIn
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Photos.Photo, nil
}

func flickrPhotoUrl(photo FlickrPhotoIn, suffix string) string {
	return fmt.Sprintf("https://farm%d.staticflickr.com/%s/%s_%s%s.jpg", photo.Farm, photo.Server, photo.Id, photo.Secret, suffix)
}

func handleFlickrSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.FormValue("q")
	photos, err := flickrSearch(query)
	if err != nil {
		jsonError(w, err)
		return
	}
	var photosOut []FlickrPhotoOut
	for _, photo := range photos {
		p := FlickrPhotoOut{
			Id:           photo.Id,
			Owner:        photo.Owner,
			Title:        photo.Title,
			SquareUrl:    flickrPhotoUrl(photo, "_q"),
			ThumbnailUrl: flickrPhotoUrl(photo, "_t"),
			MediumUrl:    flickrPhotoUrl(photo, "_z"),
			LargeUrl:     flickrPhotoUrl(photo, "_b"),
		}
		photosOut = append(photosOut, p)
	}
	b, err := json.Marshal(FlickrOk{Status: "ok", Photos: photosOut})
	if err != nil {
		jsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

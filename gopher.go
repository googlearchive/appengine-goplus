package goplus

import (
	"net/http"
	// Import google api go client library
	"code.google.com/p/google-api-go-client/googleapi/transport"
	// Import Google+ package, the package will be named "plus"
	"code.google.com/p/google-api-go-client/plus/v1"
	// Import appengine urlfetch package, that is needed to make http call to the api
	"appengine"
	"appengine/urlfetch"
)

// gopherFallback is the official gopher URL (in case we don't find any in the Google+ stream)
var gopherFallback = "http://golang.org/doc/gopher/gophercolor.png"

// init is called before the application start
func init() {
	// Register a handler for /gopher URLs.
	http.HandleFunc("/gopher", gopher)
}

// gopher is an HTTP handler that searches Google+ for an activity
// with a Gopher photo and redirects to the image thumbnail.
func gopher(w http.ResponseWriter, r *http.Request) {
	// Create appengine context, needed to do urlfetch call
	c := appengine.NewContext(r)

	// Create a new http client, supplying the API key we
	// generated to identify our application, and the urlfetch
	// transport necessary to make HTTP calls on App Engine
	transport := &transport.APIKey{
		Key:       apiKey,
		Transport: &urlfetch.Transport{Context: c}}
	client := &http.Client{Transport: transport}

	// Create the plus service object with which we can make an API call
	service, err := plus.New(client)
	if err != nil {
		// error creating the Google+ client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Search the stream for "gopher" related activities
	result, err := service.Activities.Search("#gopher").Do()
	if err != nil {
		// error searching Google+ for gopher
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Iterate through the activities search result until we find a photo or album attachment
	for _, item := range result.Items {
		for _, att := range item.Object.Attachments {
			switch att.ObjectType {
			case "photo":
				// Redirect to the gopher thumbnail
				http.Redirect(w, r, att.Image.Url, 302)
				return
			case "album":
				http.Redirect(w, r, att.Thumbnails[0].Image.Url, 302)
				return
			}
		}
	}
	// Fallback on the official gopher if we didn't found any gophers in the stream
	http.Redirect(w, r, gopherFallback, http.StatusFound)
}
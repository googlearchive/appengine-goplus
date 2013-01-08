package goplus

import (
	"fmt"
	"net/http"
)

// init is called before the application starts.
func init() {
	// Register a handler for /hello URLs.
	http.HandleFunc("/", hello)
}

// hello is an HTTP handler that prints "Hello Gopher!"
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Gopher!")
}

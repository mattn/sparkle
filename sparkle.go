package sparkle

import (
	"net/http"
)

// ListenAndServe sets up a http socket on a given address and starts
// listening for incoming requests
func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}

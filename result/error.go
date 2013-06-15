package result

import (
	"github.com/sekhat/sparkle"
	"net/http"
)

type errorResult struct {
	message string
	code    int
}

// Execute writes the error described in the errorResult to the http.ResponseWriter
func (res *errorResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	http.Error(w, res.message, res.code)
	return nil
}

// Error returns an ActionResult that when executed will write an error to the
// http ResponseWriter
func Error(message string, code int) sparkle.ActionResult {
	return &errorResult{message, code}
}

package result

import (
	"github.com/sekhat/sparkle"
	"net/http"
)

type redirectResult struct {
	url  string
	code int
}

// Redirect returns an ActionResult that sends a redirection to the client
func Redirect(url string, code int) sparkle.ActionResult {
	return &redirectResult{url, code}
}

// Execute performs the redirect described in the redirectResult
func (res *redirectResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	http.Redirect(w, r, res.url, res.code)
	return nil
}

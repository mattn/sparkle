package results

import (
	"net/http"
	"sparkle"
)

type redirectResult struct {
	url string
	code int
}

func Redirect(url string, code int) sparkle.ActionResult {
	return &redirectResult{url, code}
}

func (res *redirectResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	http.Redirect(w, r, res.url, res.code)
	return nil
}
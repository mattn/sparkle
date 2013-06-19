package result

import (
	"encoding/json"
	"github.com/sekhat/sparkle"
	"net/http"
)

type jsonResult struct {
	data interface{}
}

// JSON returns an ActionResult that sends json to the client based on a supplied model
func JSON(data interface{}) sparkle.ActionResult {
	return &jsonResult{data}
}

// Execute performs the redirect described in the redirectResult
func (res *jsonResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	w.Header().Add("content-type", "text/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	return enc.Encode(res.data)
}

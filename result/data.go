package result

import (
	"bytes"
	"github.com/sekhat/sparkle"
	"io"
	"net/http"
)

type dataResult struct {
	contentType string
	data        io.Reader
}

// Data returns an ActionResult that sends data to the client, with the specified 
// Content-Type header value
func Data(contentType string, data io.Reader) sparkle.ActionResult {
	return &dataResult{contentType, data}
}

// Bytes returns an ActionResult that sends data to the client, with the specified 
// Content-Type header value
func Bytes(contentType string, data []byte) sparkle.ActionResult {
	return Data(contentType, bytes.NewReader(data))
}

func (res *dataResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	w.Header().Add("content-type", res.contentType)
	_, err := io.Copy(w, res.data)
	return err
}

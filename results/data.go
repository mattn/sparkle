package results

import (
	"bytes"
	"io"
	"net/http"
	"github.com/sekhat/sparkle"
)

type dataResult struct {
	contentType string
	data        io.Reader
}

// Returns an ActionResult that sends data to the client, with the specified 
// Content-Type header value
func Data(contentType string, data io.Reader) sparkle.ActionResult {
	return &dataResult{contentType, data}
}

// Returns an ActionResult that sends data to the client, with the specified 
// Content-Type header value
//
func Bytes(contentType string, data []byte) sparkle.ActionResult {
	return Data(contentType, bytes.NewReader(data))
}

func (res *dataResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	w.Header().Add("content-type", res.contentType)

	buffer := make([]byte, 4096)

	for true {
		read, err := res.data.Read(buffer)

		if read > 0 {
			if _, err := w.Write(buffer[:read]); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

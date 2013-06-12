package results

import(
	"net/http"
	"github.com/sekhat/sparkle"
)

type errorResult struct {
	message string
	code    int
}

func (res *errorResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	http.Error(w, res.message, res.code)
	return nil
}

func Error(message string, code int) sparkle.ActionResult {
	return &errorResult{message, code}
}
package view

import (
	"github.com/sekhat/sparkle"
	"net/http"
)

type viewResult struct {
	viewName string
	model    interface{}
}

// View returns an Action result for rendering a ViewWriter.
func View(name string, model interface{}) sparkle.ActionResult {
	return &viewResult{
		name,
		model,
	}
}

func (res *viewResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	view, err := Get(res.viewName)
	if err != nil {
		return err
	}

	// This is probably a bad assumption that all views will return html
	// I'll come back and move this else where later
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	return view.Execute(w, res.model)
}

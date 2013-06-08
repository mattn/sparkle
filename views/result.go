package views

import (
	"errors"
	"net/http"
	"sparkle"
)

type viewResult struct {
	viewName string
	model    interface{}
}

func (res *viewResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	view, ok := registeredViews[res.viewName]
	if !ok {
		return errors.New("No view was registered with that name")
	}

	// This is probably a bad assumption that all views will return html
	// I'll come back and move this else where later
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	return view.Execute(w, res.model)
}

// Returns an Action result for rendering a View.
//
// The name is the name of the view writer to use, and model is
// the type the view writer will use for data when writing it's
// view
func View(name string, model interface{}) sparkle.ActionResult {
	return &viewResult{
		name,
		model,
	}
}

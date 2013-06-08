package views

import (
	"net/http"
	"sparkle"
	"errors"
)

type viewResult struct {
	viewName string
	model interface{} 
}

func (res *viewResult) Execute(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	view, ok := registeredViews[res.viewName]
	if !ok {
		return errors.New("No view was registered with that name")
	}	

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	return view.Execute(w, res.model)
}

func View(name string, model interface{}) sparkle.ActionResult {
	return &viewResult{
		name, 
		model,
	}
}
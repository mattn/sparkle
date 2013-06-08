package views

import (
	"errors"
	"io"

)

type ViewWriter interface {	
	Execute(io.Writer, interface{}) error
}

var registeredViews map[string]ViewWriter;

func init() {
	registeredViews = make(map[string]ViewWriter)
}

func Register(viewName string, view ViewWriter) error {
	if _, ok := registeredViews[viewName]; ok {
		return errors.New("A view with that name has already been registered")
	}

	registeredViews[viewName] = view
	return nil
}
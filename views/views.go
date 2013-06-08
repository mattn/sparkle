package views

import (
	"errors"
	"io"
)

type ViewWriter interface {
	Execute(io.Writer, interface{}) error
}

var registeredViews map[string]ViewWriter

func init() {
	registeredViews = make(map[string]ViewWriter)
}

// Register a ViewWriter under a given name
//
// This function returns an error if a ViewWriter has already been registered
// under the given name
//
func Register(viewName string, view ViewWriter) error {
	if _, ok := registeredViews[viewName]; ok {
		return errors.New("A ViewWriter with that name has already been registered")
	}

	registeredViews[viewName] = view
	return nil
}

func Get(viewName string) (ViewWriter, error) {
	view, ok := registeredViews[viewName]

	if !ok {
		return nil, errors.New("No ViewWriter with that name has been registered")
	}

	return view
}

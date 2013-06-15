package view

import (
	"errors"
	"io"
)

var NameTaken = errors.New("The name has already been registered.")

var InvalidName = errors.New("Invalid name supplied.")

type ViewWriter interface {
	Execute(io.Writer, interface{}) error
}

var registeredViews map[string]ViewWriter

func init() {
	registeredViews = make(map[string]ViewWriter)
}

// Register registers a ViewWriter under a given name
//
// The error NameTaken is returned if the supplied viewName has already
// been used to register another ViewWriter
func Register(viewName string, view ViewWriter) error {
	if _, ok := registeredViews[viewName]; ok {
		return NameTaken
	}

	registeredViews[viewName] = view
	return nil
}

// Get returns the ViewWriter
//
// The error InvalidName is returned if the viewName supplied has not
// been registered to a ViewWriter
func Get(viewName string) (ViewWriter, error) {
	view, ok := registeredViews[viewName]

	if !ok {
		return nil, InvalidName
	}

	return view, nil
}

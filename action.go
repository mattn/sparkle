package sparkle

import (
	"net/http"
)

type ActionResult interface {
	Execute(http.ResponseWriter, *http.Request, *Context) error
}

type ActionHandler func(*Context) (ActionResult, error)
type ActionWrapper func(*Context, ActionHandler) (ActionResult, error)

// Adds a request handler. By default, the pattern and matching
// is the same as the DefaultMux used by net/http
func Action(pattern string, handler ActionHandler) {
	http.HandleFunc(pattern, createActionHandler(handler))
}

func applyActionWrapper(handler ActionHandler, wrapper ActionWrapper) ActionHandler {
	return func(c *Context) (ActionResult, error) {
		return wrapper(c, handler)
	}
}

// Applies the Request Wrappers to the given Request Handler
func ApplyActionWrappers(handler ActionHandler, wrappers ...ActionWrapper) ActionHandler {
	result := handler

	for _, wrapper := range wrappers {
		result = applyActionWrapper(handler, wrapper)
	}

	return result
}

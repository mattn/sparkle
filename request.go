package sparkle

import (
	"net/http"
	"errors"
)

type RequestInitHookFunc func(*Context) error

var requestInitHooks []RequestInitHookFunc

func init() {
	requestInitHooks = make([]RequestInitHookFunc, 0)
}

// AddRequestInitHook adds a func to the list of request initializers
// 
// All request initializer will be run immediately after creation of the
// request context.
func AddRequestInitHook(hook RequestInitHookFunc) {
	requestInitHooks = append(requestInitHooks, hook)
}

func callErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func callModuleRequestInitHooks(c *Context) error {
	for _, v := range requestInitHooks {
		err := v(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func createAndInitializeContext(w http.ResponseWriter, r *http.Request) *Context {
	c := newContext(w, r)
	callModuleRequestInitHooks(c)
	return c
}

func callActionHandler(h ActionHandler, c *Context) (ActionResult, error) {
	result, err := h(c);

	if err != nil {
		return nil, err
	} 

	if result == nil {
		return nil, errors.New("No result returned from ActionHandler")
	}

	return result, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request, h ActionHandler) {
	c := createAndInitializeContext(w, r)

	result, err := callActionHandler(h, c)
	if err != nil {
		callErrorHandler(w, r, err)
		return
	}

	if err := result.Execute(w, r, c); err != nil {
		callErrorHandler(w, r, err)
		return		
	}
}
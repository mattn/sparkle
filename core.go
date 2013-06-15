package sparkle

import (
	"errors"
	"net/http"
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

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
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

func handleRequest(w http.ResponseWriter, r *http.Request, h ActionHandler) error {
	c := newContext(w, r)
	callModuleRequestInitHooks(c)

	result, err := h(c)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("No result returned from ActionHandler")
	}

	return result.Execute(w, r, c)
}

func createActionHandler(h ActionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r, h); err != nil {
			callErrorHandler(w, r, err)
		}
	}
}

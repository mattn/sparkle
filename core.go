package sparkle

import (
	"errors"
	"fmt"
	"net/http"
)

type ActionResult interface {
	Execute(http.ResponseWriter, *http.Request, *Context) error
}

type RequestHandler func(*Context) (ActionResult, error)

type RequestInitHookFunc func(http.ResponseWriter, *http.Request, *Context) error

var requestInitHooks []RequestInitHookFunc

func init() {
	requestInitHooks = make([]RequestInitHookFunc, 0)
}

func AddRequestInitHook(hook RequestInitHookFunc) {
	requestInitHooks = append(requestInitHooks, hook)
}
func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func AddHandler(path string, handler RequestHandler) {
	http.HandleFunc(path, createRequestHandler(handler))
}

func callErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, fmt.Sprint(err), 503)
}

func callModuleRequestInitHooks(w http.ResponseWriter, r *http.Request, c *Context) error {
	for _, v := range requestInitHooks {
		err := v(w, r, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request, h RequestHandler) error {
	c := newContext()
	result, err := h(c)

	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("No result returned from RequestHandler")
	}

	return result.Execute(w, r, c)
}

func createRequestHandler(h RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r, h); err != nil {
			callErrorHandler(w, r, err)
		}
	}
}

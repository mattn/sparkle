package sparkle

import (
	"errors"	
	"net/http"
)

type ActionResult interface {
	Execute(http.ResponseWriter, *http.Request, *Context) error
}

// A Function definition for Request Handlers
type RequestHandler func(*Context)(ActionResult, error)
// A Function definition for Request Handler Wrappers
//
// A RequestHandlerWrapper is used for 
type RequestHandlerWrapper func(*Context, RequestHandler)(ActionResult, error)
// 
type RequestInitHookFunc func(*Context) error

var requestInitHooks []RequestInitHookFunc

func init() {
	requestInitHooks = make([]RequestInitHookFunc, 0)
}

// Add a Request Initialization Hook
//
// Request initialization hooks are called at the start of a request
//
// See sparkle/auth as it uses this to retrieve and check and set authentication
// information before the actual handling of the request begins
func AddRequestInitHook(hook RequestInitHookFunc) {
	requestInitHooks = append(requestInitHooks, hook)
}


func applyRequestWrapper(handler RequestHandler, wrapper RequestHandlerWrapper) RequestHandler {
	return func(c *Context)(ActionResult, error) {
		return wrapper(c, handler)
	}
}

// Applies the Request Wrappers to the given Request Handler
func ApplyRequestWrappers(handler RequestHandler, wrappers ...RequestHandlerWrapper) RequestHandler {
	result := handler

	for _, wrapper := range wrappers {
		result = applyRequestWrapper(handler, wrapper)
	}

	return result
}

// Begins listening for http requests at addr, and hands them
// to sparkle
func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}

// Adds a request handler. By default, the pattern and matching
// is the same as the DefaultMux used by net/http
func AddHandler(pattern string, handler RequestHandler) {
	http.HandleFunc(pattern, createRequestHandler(handler))
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

func handleRequest(w http.ResponseWriter, r *http.Request, h RequestHandler) error {
	c := newContext(w, r)
	callModuleRequestInitHooks(c)
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

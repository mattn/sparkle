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

// Add a Request Initialization Hook
//
// Request initialization hooks are called at the start of a request
//
// See sparkle/auth as it uses this to retrieve and check and set authentication
// information before the actual handling of the request begins
func AddRequestInitHook(hook RequestInitHookFunc) {
	requestInitHooks = append(requestInitHooks, hook)
}




// Begins listening for http requests at addr, and hands them
// to sparkle
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

func createRequestHandler(h ActionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r, h); err != nil {
			callErrorHandler(w, r, err)
		}
	}
}

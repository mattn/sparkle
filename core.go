package sparkle

import ( 
	"net/http"
	"fmt"
	"errors"
)

type ActionResult interface {
	Execute(http.ResponseWriter, *http.Request, *Context) error
}

type RequestHandler func(*Context) (ActionResult, error)
type RequestInitHookFunc func(http.ResponseWriter, *http.Request, *Context) error

var requestInitHooks []RequestInitHookFunc
var globalActionHooks []ActionHook

func init() {
	requestInitHooks = make([]RequestInitHookFunc, 0)
	globalActionHooks = make([]ActionHook, 0)
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
	for _, v := range(requestInitHooks) {
		err := v(w, r, c)
		if (err != nil) {
			return err
		}
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.RequestHandler, h RequestHandler) error {
	c := newContext()
	
	if result, err := h(c); err != nil {
		return err;
	}

	if err := result.Execute(w, r, c); err != nil {
		return err;
	}
}

func createRequestHandler(h RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r, h); err != nil {
			callErrorHandler(w, r, err);
		}
	}
}


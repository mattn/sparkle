package sparkle

import (
	"net/http"
)
// ActionResults are the results returned from ActionHandlers,
// they are used to abstract away the commonalities around writting
// responses to the client
type ActionResult interface {
	Execute(http.ResponseWriter, *http.Request, *Context) error
}
// ActionHandlers are functions that are called to handle a request in
// sparkle. They recieve a pointer to a sparkle.Context and can return
// an ActionResult and/or an error.
//
// When ActionHandlers return ActionResults (with no error), that ActionResult
// is later executed (in most cases) in order to write output to the client
//
// On errors, the framework executes it's configured error handler.
type ActionHandler func(*Context) (ActionResult, error)
// ActionWrappers are functions used to wrap ActionHandlers in order to perform
// some function prior to or after ActionHandlers.
//
// They are given a pointer to a *sparkle.Context and the ActionHandler it's wrapping.
//
// Because the ActionWrapper is responsible for calling the ActionHandler, it's possible
// for the ActionWrapper to actually prevent the ActionHandler from being run.
// Since the ActionWrapper also returns an ActionResult, it's also possible for an
// ActionWrapper to change or subvert the ActionResult to be used.
type ActionWrapper func(*Context, ActionHandler) (ActionResult, error)
// Action adds an ActionHandler to be executed when incoming urls match the given pattern.
//
// By default, the pattern matching is the same as the DefaultMux used by net/http
func Action(pattern string, handler ActionHandler) {
	http.HandleFunc(pattern, createActionHandler(handler))
}

func applyActionWrapper(handler ActionHandler, wrapper ActionWrapper) ActionHandler {
	return func(c *Context) (ActionResult, error) {
		return wrapper(c, handler)
	}
}

// ApplyActionWrappers applies one or more ActionWrappers to an ActionHandler and returns
// another ActionHandler that will make sure the wrappers are called into.
//
// Providing that all ActionWrappers call their supplied ActionHandler, then calling
//     ApplyActionWrappers(handler, w1, w2, w3)
// can be seen as creating an ActionHandler that will cause w3 to called, then w2, then w1 and finally handler
//
// A simple example of an ActionWrapper might be for logging when a certain Action is called.
//
//     func LogAction(c *Context, next ActionHandler) (ActionResult, error) {
//	       log("Action Starting")
//         res, err := next(c)
//         log("Action Complete")
//         return res, err
//     }
//
// And this action wrapper can be applied like so:
//     func main() {
//	       sparkle.Action("/MyUrlPattern", sparkle.ApplyActionWrappers(MyHandler, LogAction))
//     }
func ApplyActionWrappers(handler ActionHandler, wrappers ...ActionWrapper) ActionHandler {
	result := handler

	for _, wrapper := range wrappers {
		result = applyActionWrapper(handler, wrapper)
	}

	return result
}

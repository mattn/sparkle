package sparkle

import (
	"net/http"
)

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	data           map[string]interface{}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w,
		r,
		make(map[string]interface{}),
	}
}

// Gets the response writer associated with the Context
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

// Gets the request associated with the Context
func (c *Context) Request() *http.Request {
	return c.request
}

// Sets a value against a Context with the given key
//
// This is primarily there so that extensions can store data
// against a context
func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

// Gets a value stored against a Context by it's key
//
// If the key does not exist, nil is returned
func (c *Context) Get(key string) interface{} {
	result, ok := c.data[key]
	if !ok {
		return nil
	}
	return result
}

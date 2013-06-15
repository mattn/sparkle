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

// ResponseWriter returns the http.ResponseWriter associated with the Context
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

// Request returns the http.Request associated with the Context
func (c *Context) Request() *http.Request {
	return c.request
}

// Set sets a value against a Context with a given key
//
// This is primarily there so that extensions can store data
// against a context
func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

// Get returns the value stored against a Context by it's key
//
// If the key does not exist, nil is returned
func (c *Context) Get(key string) interface{} {
	result, ok := c.data[key]
	if !ok {
		return nil
	}
	return result
}

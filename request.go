package sparkle

import (
	"net/http"
)

const (
	contextRequestKey        string = "Sparkle.Request"
	contextResponseWriterKey string = "Sparkle.ResponseWriter"
)

func init() {
	AddRequestInitHook(moduleRequestInit)
}

func moduleRequestInit(w http.ResponseWriter, r *http.Request, c *Context) error {
	c.Set(contextRequestKey, r)
	c.Set(contextResponseWriterKey, w)
	return nil
}

// Gets the net/http Request associated with the context
func (c *Context) Request() *http.Request {
	o := c.Get(contextRequestKey)
	r, ok := o.(*http.Request)
	if !ok {
		return nil
	}
	return r
}

// Gets the net/http ResponseWriter associated with the context
func (c *Context) ResponseWriter() http.ResponseWriter {
	o := c.Get(contextResponseWriterKey)
	r, ok := o.(http.ResponseWriter)
	if !ok {
		return nil
	}
	return r
}

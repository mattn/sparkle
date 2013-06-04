package auth

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"encoding/gob"
	"sparkle"	
	"strings"
)

const (
	authDataKey string = "Sparkle.Auth"
	authCookieName string = "Sparkle.Auth.Cookie"
)

type authData struct {
	UserIdentifier string
	Addr string
}

var sc *securecookie.SecureCookie

func init() {
	gob.Register(&authData{})
	sc = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
}

func AuthInit() {	
	sparkle.AddRequestInitHook(authInitRequestHook)
}

func authInitRequestHook(w http.ResponseWriter, r *http.Request, c *sparkle.Context) error {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		return nil
	}

	value := &authData{}
	if err := sc.Decode(authCookieName, cookie.Value, &value); err != nil {
		// Despite decoding failing, we'll assume this just means we're not authenticated
		// So no need to return an error
		return nil
	}

	c.Set(authDataKey, value)
}

func getUserIP(r *http.Request) string {
	var portSeperatorIndex = strings.Index(r.RemoteAddr, ":")	
	if (portSeperatorIndex == -1) {
		return r.RemoteAddr
	} else {
		return r.RemoteAddr[:portSeperatorIndex]
	}
}

func (a *authData) isValid(r *http.Request) bool {
	return getUserIP(r) == a.Addr &&
		   a.UserIdentifier != ""
}

func (c *sparkle.Context) getAuthData() *authData {
	data = c.Get(authDataKey)

	if result, ok := data.(authData); !ok {
		return nil
	}

	return result;
}

// Returns a boolean indictating whether the current context can be considered
// authenticated
func IsAuthenticated(c *sparkle.Context) bool {
	if auth := c.Get(authDataKey); auth == nil {
		return false
	}

	return auth.isValid()
}

// Get the user identifier of the authenticated user, or nil if there is no 
// authenticated user
func AuthenticatedAs(c *sparkle.Context) string {
	if auth := c.Get(authDataKey); auth == nil {
		return nil
	}

	if !auth.isValid() {
		return nil
	}

	return auth.UserIdentifier
}

// Sets an Authentication Cookie for a given path, and userIdentifier
func Authenticate(c *sparkle.Context, path string, userIdentifier string) error {
	r := c.Request()
	w := c.ResponseWriter()

	value := &authData{ userIdentifier, r.getUserIP() }
	if encoded, err := s.Encode(authCookieName, value); err != nil {
		return error
	}

	http.SetCookie(w, &http.Cookie{
		Name: authCookieName,
		Value: encoded,
		Path: path,
	})

	c.Set(authDataKey, auth)
}

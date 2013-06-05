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
	return nil
}

func getUserIP(r *http.Request) string {
	var portSeperatorIndex = strings.Index(r.RemoteAddr, ":")	
	if (portSeperatorIndex == -1) {
		return r.RemoteAddr
	}

	return r.RemoteAddr[:portSeperatorIndex]	
}

func (a *authData) isValid(r *http.Request) bool {
	return getUserIP(r) == a.Addr &&
		   a.UserIdentifier != ""
}

func getAuthData(c *sparkle.Context) *authData {
	data := c.Get(authDataKey)

	result, ok := data.(*authData);
	if !ok {
		return nil
	}

	return result
}

// Returns a boolean indictating whether the current context can be considered
// authenticated
func IsAuthenticated(c *sparkle.Context) bool {
	auth := getAuthData(c)
	if auth == nil {
		return false
	}

	return auth.isValid(c.Request())
}

// Get the user identifier of the authenticated user, or nil if there is no 
// authenticated user
func AuthenticatedAs(c *sparkle.Context) string {
	auth := getAuthData(c)

	if auth == nil || !auth.isValid(c.Request()) {
		return ""
	}

	return auth.UserIdentifier
}

// Sets an Authentication Cookie for a given path, and userIdentifier
func Authenticate(c *sparkle.Context, path string, userIdentifier string) error {
	r := c.Request()
	w := c.ResponseWriter()

	value := &authData{ userIdentifier, getUserIP(r) }
	encoded, err := sc.Encode(authCookieName, value) 
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name: authCookieName,
		Value: encoded,
		Path: path,
	})

	c.Set(authDataKey, value)
	return nil
}

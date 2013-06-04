package auth

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"sparkle"
	"errors"
)

const (
	authDataKey string = "Sparkle.Auth"
	authCookieName string = "Sparkle.Auth.Cookie"
)

type authData struct {
	UserIdentifier string
	Address string
}

var sc *securecookie.SecureCookie

func init() {	
	sc = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
}

func AuthInit() {	
	sparkle.AddRequestInitHook(authInitRequestHook)
}

func authInitRequestHook(w http.ResponseWriter, r *http.Request, c *sparkle.Context) {
	cookie, err := r.Cookie(authCookieName)

	// Get Cookie from request
	// if Cookie exists decode data and check validity, store Authenticated AuthIdentity against context
	// if doesn't exist store empty non authenticated AuthIdentity
}

func (c *sparkle.Context) getAuthData() *authData {
	data = c.Get(authDataKey)

	if result, ok := data.(authData); !ok {
		return nil
	}

	return result;
}

func (c *sparkle.Context) IsAuthenticated() bool {
	if auth := c.Get(authDataKey); auth == nil {
		return false;
	}

	auth.

	return false
}

func (c *sparkle.Context) AuthenticatedAs() string {
	return nil
}

// Sets a user as authenticated as a given user identifier
func (c *sparkle.Context) Authenticate(userIdentifier string) {

}

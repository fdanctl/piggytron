// Package shared provides small HTTP helpers used across handlers, such as
// CookieMaker for creating and revoking cookies.
package shared

import (
	"net/http"
	"time"
)

// CookieMaker stamps session cookies from a shared set of default values.
type CookieMaker struct {
	defaultCookie http.Cookie
}

// NewCookieMaker builds the maker from the default cookie attributes.
func NewCookieMaker(defaultValues http.Cookie) *CookieMaker {
	return &CookieMaker{defaultCookie: defaultValues}
}

// NewCookie returns a copy of the default cookie holding the given session
// value.
func (cm *CookieMaker) NewCookie(value string) *http.Cookie {
	nc := cm.defaultCookie
	nc.Value = value
	return &nc
}

// RevokeCookie returns a copy of the default cookie expired in the past,
// which the browser treats as a deletion.
func (cm *CookieMaker) RevokeCookie() *http.Cookie {
	nc := cm.defaultCookie
	nc.Expires = time.Unix(0, 0)
	nc.MaxAge = -1
	return &nc
}

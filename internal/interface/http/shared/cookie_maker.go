package shared

import (
	"net/http"
	"time"
)

type CookieMaker struct {
	defaultCookie http.Cookie
}

func NewCookieMaker(defaultValues http.Cookie) *CookieMaker {
	return &CookieMaker{defaultCookie: defaultValues}
}

func (cm *CookieMaker) NewCookie(value string) *http.Cookie {
	nc := cm.defaultCookie
	nc.Value = value
	return &nc
}

func (cm *CookieMaker) RevokeCookie() *http.Cookie {
	nc := cm.defaultCookie
	nc.Expires = time.Unix(0, 0)
	nc.MaxAge = -1
	return &nc
}

package sessions

import (
	"fmt"

	"github.com/ayushbpl10/oauth2_proxy/pkg/apis/options"
	"github.com/ayushbpl10/oauth2_proxy/pkg/apis/sessions"
	"github.com/ayushbpl10/oauth2_proxy/pkg/sessions/cookie"
)

// NewSessionStore creates a SessionStore from the provided configuration
func NewSessionStore(opts *options.SessionOptions, cookieOpts *options.CookieOptions) (sessions.SessionStore, error) {
	switch opts.Type {
	case options.CookieSessionStoreType:
		return cookie.NewCookieSessionStore(opts, cookieOpts)
	default:
		return nil, fmt.Errorf("unknown session store type '%s'", opts.Type)
	}
}

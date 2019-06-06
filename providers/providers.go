package providers

import (
	"github.com/ayushbpl10/oauth2_proxy/cookie"
	"github.com/ayushbpl10/oauth2_proxy/pkg/apis/sessions"
)

// Provider represents an upstream identity provider implementation
type Provider interface {
	Data() *ProviderData
	GetEmailAddress(*sessions.SessionState) (string, error)
	GetUserName(*sessions.SessionState) (string, error)
	Redeem(string, string) (*sessions.SessionState, error)
	ValidateGroup(string) bool
	ValidateSessionState(*sessions.SessionState) bool
	GetLoginURL(redirectURI, finalRedirect string) string
	RefreshSessionIfNeeded(*sessions.SessionState) (bool, error)
	SessionFromCookie(string, *cookie.Cipher) (*sessions.SessionState, error)
	CookieForSession(*sessions.SessionState, *cookie.Cipher) (string, error)
}

// New provides a new Provider based on the configured provider string
func New(p ProviderInput) Provider {
	switch p.ProviderName {
	case "linkedin":
		return NewLinkedInProvider(p.P)
	case "facebook":
		return NewFacebookProvider(p.P)
	case "github":
		return NewGitHubProvider(p.P)
	case "azure":
		return NewAzureProvider(p.P)
	case "gitlab":
		return NewGitLabProvider(p.P)
	case "oidc":
		return NewOIDCProvider(p.P)
	case "login.gov":
		return NewLoginGovProvider(p.P)
	default:
		return NewGoogleProvider(p.P)
	}
}

type ProviderInput struct {
	ProviderName string
	P *ProviderData
}
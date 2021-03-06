package test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/url"

	"github.com/keratin/authn-server/app"
	"github.com/keratin/authn-server/app/data/mock"
	"github.com/keratin/authn-server/lib/oauth"
	"github.com/keratin/authn-server/lib/route"
	"github.com/keratin/authn-server/ops"
)

func App() *app.App {
	authnURL, err := url.Parse("https://authn.example.com")
	if err != nil {
		panic(err)
	}

	weakKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}

	cfg := app.Config{
		BcryptCost:              4,
		SessionSigningKey:       []byte("TestKey"),
		AuthNURL:                authnURL,
		SessionCookieName:       "authn",
		OAuthCookieName:         "authn-oauth-nonce",
		ApplicationDomains:      []route.Domain{{Hostname: "test.com"}},
		PasswordMinComplexity:   2,
		AppPasswordResetURL:     &url.URL{Scheme: "https", Host: "app.example.com"},
		AppPasswordlessTokenURL: &url.URL{Scheme: "https", Host: "app.example.com"},
		EnableSignup:            true,
	}

	return &app.App{
		Config:            &cfg,
		KeyStore:          mock.NewKeyStore(weakKey),
		AccountStore:      mock.NewAccountStore(),
		RefreshTokenStore: mock.NewRefreshTokenStore(),
		Actives:           mock.NewActives(),
		Reporter:          &ops.LogReporter{},
		OauthProviders:    map[string]oauth.Provider{},
	}
}

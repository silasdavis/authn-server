package handlers_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"testing"

	"github.com/keratin/authn-server/server/test"
	"github.com/keratin/authn-server/app"
	"github.com/keratin/authn-server/app/data/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJWKs(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 512)
	require.NoError(t, err)
	app := &app.App{
		KeyStore: mock.NewKeyStore(rsaKey),
		Config:   &app.Config{},
	}

	server := test.Server(app)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/jwks", server.URL))
	require.NoError(t, err)
	body := test.ReadBody(res)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []string{"application/json"}, res.Header["Content-Type"])
	assert.NotEmpty(t, body)
}

func BenchmarkGetJWKs(b *testing.B) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	app := &app.App{
		KeyStore: mock.NewKeyStore(rsaKey),
		Config:   &app.Config{},
	}

	server := test.Server(app)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		http.Get(fmt.Sprintf("%s/jwks", server.URL))
	}
}

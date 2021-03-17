package web

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	keyUsername = "username"
	jwtCookie   = "jwt-auth"
)

type Token struct {
	Username string
}

type AuthContext struct {
	LoggedIn bool
	Username string
}

type Auth struct {
	key *rsa.PrivateKey
}

func NewAuth() (*Auth, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.WithError(err).Errorf("failed to generate private key")
		return nil, err
	}

	return &Auth{
		key: key,
	}, nil
}

func (a *Auth) IssueToken(username string) (string, error) {
	t := jwt.New()
	t.Set(keyUsername, username)

	// Signing a token (using raw rsa.PrivateKey)
	signed, err := jwt.Sign(t, jwa.RS256, a.key)
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func (a *Auth) ParseToken(payload string) (*Token, error) {
	jt, err := jwt.ParseString(
		payload,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &a.key.PublicKey),
	)
	if err != nil {
		return nil, err
	}

	username, ok := jt.Get(keyUsername)
	if !ok {
		return nil, errors.New("username not found")
	}

	ustr, ok := username.(string)
	if !ok {
		return nil, errors.New("username is not a string")
	}

	log.WithField("username", ustr).Debug("authorized request")

	token := Token{
		Username: ustr,
	}

	return &token, nil
}

func (a *Auth) FromRequest(r *http.Request) AuthContext {
	cookie, err := r.Cookie(jwtCookie)
	if err != nil {
		return AuthContext{}
	}

	tokenStr := cookie.Value
	t, err := a.ParseToken(tokenStr)
	if err != nil {
		return AuthContext{}
	}

	return AuthContext{
		LoggedIn: true,
		Username: t.Username,
	}
}

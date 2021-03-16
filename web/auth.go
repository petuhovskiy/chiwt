package web

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

const (
	keyUsername = "username"
)

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

func (a *Auth) ParseToken(payload string) (jwt.Token, error) {
	token, err := jwt.ParseString(
		payload,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &a.key.PublicKey),
	)
	if err != nil {
		return nil, err
	}

	log.WithField("token", token).Info("parsed token")

	return token, nil
}

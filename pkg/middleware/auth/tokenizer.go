package auth

import (
	"context"
	"errors"
	"time"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/golang-jwt/jwt"
)

var ErrInvalidCred error = errors.New("Invalid credentials")
var ErrInvalidToken error = errors.New("Invalid token")

type Tokenizer interface {
	CreateToken(map[string]any, time.Duration) (string, error)
	CheckToken(context.Context, string) (*tasks.User, error)
}

type jwtTokenizer struct {
	expiresDelta time.Duration
	secretKey    string
}

func (j jwtTokenizer) CreateToken(data map[string]any, exp time.Duration) (string, error) {
	if exp == 0 {
		exp = j.expiresDelta
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.StandardClaims{
		ExpiresAt: time.Now().Add(exp).Unix(),
		Subject:   data["sub"].(string),
	}).SignedString(j.secretKey)
}

func (j jwtTokenizer) CheckToken(ctx context.Context, token string) (*tasks.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Couldn't parse claims: Not map claims")
	}

	username, ok := claims["sub"].(string)
	if !ok || username == "" {
		return nil, ErrInvalidCred
	}

	// user, err := get user from db
	var user *tasks.User

	if err != nil {
		return nil, ErrInvalidCred
	}
	return user, nil
}

func NewTokenizer(expiresDelta time.Duration, secretKey string) Tokenizer {
	return jwtTokenizer{
		expiresDelta: expiresDelta,
		secretKey:    secretKey,
	}
}

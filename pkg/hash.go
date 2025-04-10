package tasks

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	HashPassword(string) (string, error)
	CheckPassword(plain, hashed string) bool
}

type bcryptHasher struct{}

func (b bcryptHasher) HashPassword(password string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(s), err
}

func (b bcryptHasher) CheckPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func NewHasher() Hasher {
	return bcryptHasher{}
}

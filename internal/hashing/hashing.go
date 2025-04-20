package hashing

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PassHasherInterface interface {
	HashPass(password string) (string, error)
	ComparePass(password, hashedPassword string) error
}

type PassHasher struct {
}

func NewPassHasher() PassHasherInterface {
	return &PassHasher{}
}
func (p *PassHasher) HashPass(password string) (string, error) {
	HashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash pass %v", err)
	}
	return string(HashedPass), nil
}
func (p *PassHasher) ComparePass(password, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

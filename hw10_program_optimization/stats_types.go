package hw10programoptimization

import (
	"errors"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

type users [100_000]User

var (
	ErrDomainNotFound = errors.New("domain not found in user email")
	ErrInvalidEmail   = errors.New("invalid email")
)
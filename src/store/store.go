package store

import (
	"github.com/phuwn/crawlie/src/store/keyword"
	"github.com/phuwn/crawlie/src/store/user"
)

// Store - server store struct
type Store struct {
	User    user.Store
	Keyword keyword.Store
}

// New - create new store variable
func New() *Store {
	return &Store{
		User:    user.NewStore(),
		Keyword: keyword.NewStore(),
	}
}

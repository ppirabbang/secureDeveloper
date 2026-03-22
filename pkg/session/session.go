package session

import (
	"crypto/rand"
	"encoding/hex"
)

const sessionTokenBytes = 24

type Store struct {
	tokens map[string]uint
}

func NewStore() *Store {
	return &Store{
		tokens: make(map[string]uint),
	}
}

func (s *Store) Create(userID uint) (string, error) {
	token, err := newToken()
	if err != nil {
		return "", err
	}
	s.tokens[token] = userID
	return token, nil
}

func (s *Store) Lookup(token string) (uint, bool) {
	userID, ok := s.tokens[token]
	return userID, ok
}

func (s *Store) Delete(token string) {
	delete(s.tokens, token)
}

func newToken() (string, error) {
	buffer := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

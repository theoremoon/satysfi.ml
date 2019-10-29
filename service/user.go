package service

import (
	"github.com/theoremoon/SATySFi-Online/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(username, password string) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err := s.repo.InsertUser(username, string(hashed))
	if err != nil {
		return nil, err
	}
	return user, nil
}

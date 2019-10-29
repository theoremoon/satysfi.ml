package service

import (
	"math/rand"
	"time"

	"github.com/theoremoon/SATySFi-Online/model"
	"github.com/theoremoon/SATySFi-Online/repository"
)

type Service interface {
	Register(username, password string) (*model.User, error)
	CreateSession(userID int) (*model.Session, error)
	Close() error
}

func New(repo repository.Repository, seed int64) Service {
	return &service{
		rnd:  rand.New(rand.NewSource(seed)),
		repo: repo,
	}
}

type service struct {
	rnd  *rand.Rand
	repo repository.Repository
}

func (s *service) Close() error {
	return s.repo.Close()
}

func (s *service) CreateSession(userID int) (*model.Session, error) {
	expiredAt := time.Now().AddDate(0, 1, 0).Unix()
	sessionID := s.rnd.Int()

	session, err := s.repo.InsertSession(userID, sessionID, expiredAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

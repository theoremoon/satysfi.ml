package repository

import (
	"github.com/theoremoon/SATySFi-Online/model"
)

type SessionRepository interface {
	InsertSession(userID, sessionID int, expiredAt int64) (*model.Session, error)
	GetValidSessionByID(sessionID int) (*model.Session, error)
}

func (r *repository) InsertSession(userID, sessionID int, expiredAt int64) (*model.Session, error) {
	session := model.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiredAt: expiredAt,
	}

	_, err := r.db.NamedExec(`
		INSERT INTO sessions(id, userid, expired_at)
		VALUES (:id, :user_id, :expired_at)
	`, session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) GetValidSessionByID(sessionID int) (*model.Session, error) {
	var session model.Session
	err := r.db.Get(
		&session,
		`
		SELECT (id, user_id, expired_at)
		FROM sessions
		WHERE id = ? AND expired_at < extract(epoc FROM now())
		LIMIT 1
		`)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

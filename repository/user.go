package repository

import (
	"github.com/theoremoon/SATySFi-Online/model"
)

type UserRepository interface {
	InsertUser(username, hashedpassword string) (*model.User, error)
	GetUserByID(userID int) (*model.User, error)
}

func (r *repository) InsertUser(username, hashedpassword string) (*model.User, error) {
	user := model.User{
		Username: username,
		Password: hashedpassword,
	}

	rows, err := r.db.NamedQuery(`
		INSERT INTO users(username, password)
		VALUES (:username, :password)
		RETURNING id
	`, user)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(&user.ID)
	}

	return &user, nil
}

func (r *repository) GetUserByID(userID int) (*model.User, error) {
	var user model.User

	err := r.db.Get(&user, `SELECT id, username FROM users WHERE id = ? LIMIT 1`, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

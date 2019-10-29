package model

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"` // maybe password hash
}

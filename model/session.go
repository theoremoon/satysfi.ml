package model

type Session struct {
	ID        int   `db:"id"`
	UserID    int   `db:"user_id"`
	ExpiredAt int64 `db:"expired_at"`
}

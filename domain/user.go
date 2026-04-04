package domain

import "time"

type User struct {
	ID            int64      `db:"id"`
	Email         string     `db:"email"`
	PasswordHash  string     `db:"password_hash"`
	EmailVerified bool       `db:"email_verified"`
	Status        string     `db:"status"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

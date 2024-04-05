package models

import "time"

type UserModel struct {
	Id        uint32    `db:"id"         json:"id"`
	Username  string    `db:"username"   json:"username"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

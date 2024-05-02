package models

import "time"

type UserModel struct {
	Id        uint32    `db:"id"         json:"id"        example:"1"`
	Username  string    `db:"username"   json:"username"  example:"NahtanN"`
	Email     string    `db:"email"      json:"email"     example:"nahtann@outlook.com"`
	Password  string    `db:"password"   json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" example:"2024-05-02T00:10:10.875334-03:00"`
}

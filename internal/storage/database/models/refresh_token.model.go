package models

import "time"

type RefreshTokenModel struct {
	Id            string    `db:"id"`
	ParentTokenId string    `db:"parent_token_id"`
	Token         string    `db:"token"`
	UserId        string    `db:"user_id"`
	ExpiresAt     time.Time `db:"expires_at"`
	CreatedAt     time.Time `db:"created_at"`

	User UserModel
}

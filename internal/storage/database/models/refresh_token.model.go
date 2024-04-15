package models

import "time"

type RefreshTokenModel struct {
	Id            uint32    `db:"id"`
	ParentTokenId uint32    `db:"parent_token_id"`
	Token         string    `db:"token"`
	UserId        uint32    `db:"user_id"`
	ExpiresAt     time.Time `db:"expires_at"`
	Used          bool      `db:"used"`
	CreatedAt     time.Time `db:"created_at"`

	User UserModel
}

func (rt *RefreshTokenModel) Table() (string, string) {
	return "refresh_tokens", "rt"
}

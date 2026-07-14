package model

type MfaList struct {
	Totp bool `json:"totp"`
}

type MfaTotp struct {
	ID     int32  `db:"id"      json:"id"`
	Key    string `db:"key"     json:"key"`
	UserID string `db:"user_id" json:"user_id"`
}

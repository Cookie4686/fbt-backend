package model

type MfaList struct {
	Totp bool `json:"totp"`
}

type MfaTotp struct {
	ID     int32  `json:"id" db:"id"`
	Key    string `json:"key" db:"key"`
	UserID string `json:"user_id" db:"user_id"`
}

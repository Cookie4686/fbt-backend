package model

import authv1 "fbt/backend/gen/proto/go/auth/v1"

type Passkey struct {
	PasskeyID      string   `db:"passkey_id"`
	PublicKey      []byte   `db:"public_key"`
	UserID         string   `db:"user_id"`
	WebauthnUserID string   `db:"webauthn_user_id"`
	Counter        int64    `db:"counter"`
	DeviceType     string   `db:"device_type"`
	BackedUp       bool     `db:"backed_up"`
	Transports     []string `db:"transports"`
}

func (s *Passkey) ToProto() *authv1.Passkey {
	return &authv1.Passkey{
		PasskeyId:  s.PasskeyID,
		PublicKey:  s.PublicKey,
		UserId:     s.UserID,
		WebauthnId: s.WebauthnUserID,
		Counter:    s.Counter,
		DeviceType: s.DeviceType,
		BackedUp:   s.BackedUp,
		Transports: s.Transports,
	}
}

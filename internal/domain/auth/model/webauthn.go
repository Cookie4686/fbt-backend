package model

import authv1 "fbt/backend/gen/proto/go/auth/v1"

type UserWebAuthn struct {
	WebauthnID int32  `db:"webauthn_id"`
	UserID     string `db:"user_id"`
	RpID       string `db:"rp_id"`
}

type WebAuthnCredential struct {
	ID             int64    `db:"id"`
	RpID           string   `db:"rp_id"`
	UserID         string   `db:"user_id"`
	CredentialID   string   `db:"credential_id"`
	PublicKey      []byte   `db:"public_key"`
	Counter        int64    `db:"counter"`
	Aaguid         []byte   `db:"aaguid"`
	DeviceType     string   `db:"device_type"`
	Transports     []string `db:"transports"`
	UserPresent    bool     `db:"user_present"`
	UserVerified   bool     `db:"user_verified"`
	BackupEligible bool     `db:"backup_eligible"`
	BackupState    bool     `db:"backup_state"`
}

func (s *WebAuthnCredential) ToProto() *authv1.WebAuthnCredential {
	return &authv1.WebAuthnCredential{
		RpId:         s.RpID,
		UserId:       s.UserID,
		CredentialId: s.CredentialID,
		PublicKey:    s.PublicKey,
		Counter:      s.Counter,
		Aaguid:       s.Aaguid,
		DeviceType:   s.DeviceType,
		BackedUp:     s.BackupState,
		Transports:   s.Transports,
	}
}

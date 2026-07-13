-- +goose Up
CREATE TABLE webauthn_credentials (
  id               BIGSERIAL PRIMARY KEY,
  rp_id            VARCHAR(512) NOT NULL,
  user_id          VARCHAR(255) NOT NULL,
	credential_id    TEXT NOT NULL,
  public_key       BYTEA NOT NULL,
	counter          BIGINT NOT NULL,
  aaguid           BYTEA NOT NULL,
  device_type      VARCHAR(32) NOT NULL,
  transports       VARCHAR(255)[] NOT NULL,
  user_present     BOOLEAN NOT NULL DEFAULT FALSE,
  user_verified    BOOLEAN NOT NULL DEFAULT FALSE,
  backup_eligible  BOOLEAN NOT NULL DEFAULT FALSE,
  backup_state     BOOLEAN NOT NULL DEFAULT FALSE,

	CONSTRAINT "webauthn_credentials_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE,
	CONSTRAINT "webauthn_credentials_rp_id_credential_id_unique" UNIQUE("rp_id", "credential_id")
);

-- +goose Down
DROP TABLE IF EXISTS "webauthn_credentials" CASCADE;
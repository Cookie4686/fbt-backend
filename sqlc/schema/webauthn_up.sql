CREATE TABLE user_passkey (
  passkey_id       TEXT           NOT NULL,
  public_key       BYTEA          NOT NULL,
  user_id          VARCHAR(255)   NOT NULL,
  webauthn_user_id TEXT           NOT NULL,
  counter          BIGINT         NOT NULL,
  device_type      VARCHAR(32)    NOT NULL,
  backed_up        BOOLEAN        NOT NULL,
  transports       VARCHAR(255)[] NOT NULL,

	CONSTRAINT "user_passkey_passkey_id_pk" PRIMARY KEY("passkey_id"),
	CONSTRAINT "user_passkey_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE,
	CONSTRAINT "user_passkey_user_id_webauthn_user_id_unique" UNIQUE("user_id", "webauthn_user_id")
);

-- +goose Up
CREATE TABLE "mfa_totp" (
  "id" serial PRIMARY KEY,
  "key" varchar(255) NOT NULL,
	"user_id" varchar(255) NOT NULL,
	CONSTRAINT "mfa_totp_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE,
	CONSTRAINT "mfa_totp_user_id_unique" UNIQUE("user_id")
);

-- +goose Down
DROP TABLE IF EXISTS "mfa_totp" CASCADE;
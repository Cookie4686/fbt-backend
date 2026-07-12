-- +goose Up
CREATE TABLE "oauth_providers" (
    "oauth_provider_id" smallserial PRIMARY KEY,
    "name" varchar(63),
	CONSTRAINT "oauth_provider_name_unique" UNIQUE("name")
);

CREATE TABLE "user_oauth" (
	"user_id" varchar(255) NOT NULL,
  "oauth_provider_id" smallint NOT NULL,
  "id_token" varchar(255) NOT NULL,
	CONSTRAINT "user_oauth_oauth_providre_id_token_id_unique" UNIQUE( "oauth_provider_id", "id_token"),
	CONSTRAINT "user_oauth_oauth_provider_id_fk" FOREIGN KEY("oauth_provider_id") REFERENCES oauth_providers("oauth_provider_id") ON DELETE CASCADE,
	CONSTRAINT "user_oauth_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE
);

CREATE TABLE "oauth_registration" (
  "registration_id" varchar(127) PRIMARY KEY,
  "id_token" varchar(255) NOT NULL,
  "oauth_provider_id" smallint NOT NULL,
	"email_verified" boolean NOT NULL DEFAULT false,
	"expires_at" timestamp NOT NULL,
	CONSTRAINT "oauth_registration_oauth_provider_id_fk" FOREIGN KEY("oauth_provider_id") REFERENCES oauth_providers("oauth_provider_id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "oauth_providers" CASCADE;
DROP TABLE IF EXISTS "user_oauth" CASCADE;
DROP TABLE IF EXISTS "oauth_registration" CASCADE;
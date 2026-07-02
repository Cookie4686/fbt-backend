CREATE TABLE "users" (
	"user_id" varchar(255) NOT NULL,
	"username" varchar(32) NOT NULL,
	"email" varchar(255) NOT NULL,
	"email_verified" boolean NOT NULL DEFAULT false,
	"password" varchar(255),
  "password_salt" varchar(255),
	"password_enabled" boolean NOT NULL,
	CONSTRAINT "user_user_id_pk" PRIMARY KEY("user_id"),
	CONSTRAINT "user_username_unique" UNIQUE("username"),
	CONSTRAINT "user_email_unique" UNIQUE("email")
);

CREATE TABLE "sessions" (
	"session_id" varchar(255) NOT NULL,
	"user_id" varchar(255) NOT NULL,
	"expires_at" timestamp NOT NULL,
	"two_factor_verified" boolean NOT NULL DEFAULT false,
	CONSTRAINT "session_session_id_pk" PRIMARY KEY("session_id"),
	CONSTRAINT "session_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE
);

CREATE TABLE "email_verification" (
	"user_id" varchar(255) NOT NULL,
	"verification_id" varchar(255) NOT NULL,
	"otp" varchar(15) NOT NULL,
	"expires_at" timestamp NOT NULL,
	CONSTRAINT "email_verification_user_id_pk" PRIMARY KEY("user_id"),
	CONSTRAINT "email_verification_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE
);

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

CREATE TABLE "mfa_totp" (
  "id" serial PRIMARY KEY,
  "key" varchar(255) NOT NULL,
	"user_id" varchar(255) NOT NULL,
	CONSTRAINT "mfa_totp_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE,
	CONSTRAINT "mfa_totp_user_id_unique" UNIQUE("user_id")
);

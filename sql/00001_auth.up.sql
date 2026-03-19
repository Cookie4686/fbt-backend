CREATE TABLE "users" (
	"user_id" varchar(255) NOT NULL,
	"username" varchar(32) NOT NULL,
	"email" varchar(255) NOT NULL,
	"password" varchar(255),
    "password_salt" varchar(255),
	"email_verified" boolean NOT NULL DEFAULT false,
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

/* CREATE TABLE "email_verification_request" (
	"email_verification_request_id" varchar(255) NOT NULL,
	"user_id" varchar(255) NOT NULL,
	"email" varchar(255) NOT NULL,
	"code" varchar(255) NOT NULL,
	"expires_at" datetime NOT NULL,
	CONSTRAINT "email_verification_request_email_verification_request_id" PRIMARY KEY("email_verification_request_id")
) */;
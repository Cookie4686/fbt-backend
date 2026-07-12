-- +goose Up
CREATE TABLE "accounts" (
	"account_id" serial,
	"name" varchar(255) NOT NULL,
	"is_debit" boolean NOT NULL,
	"user_id" varchar(255) NOT NULL,

	CONSTRAINT "accounts_account_id_pk" PRIMARY KEY("account_id"),
	CONSTRAINT "accounts_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE
);

CREATE TABLE "tags" (
    "tag_id" bigserial,
    "name" varchar(255) NOT NULL,
	"user_id" varchar(255) NOT NULL,

	CONSTRAINT "tags_tag_id_pk" PRIMARY KEY("tag_id"),
	CONSTRAINT "accounts_user_id_fk" FOREIGN KEY("user_id") REFERENCES users("user_id") ON DELETE CASCADE,
	CONSTRAINT "accounts_name_user_id_unique" UNIQUE(name, user_id)
);

CREATE TABLE "account_tag" (
	"account_id" int NOT NULL,
    "tag_id" bigint NOT NULL,

	CONSTRAINT "account_tag_unique" UNIQUE(account_id, tag_id)
);

CREATE TABLE "transactions" (
    "transaction_id" bigserial,
    "datetime" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT "transactions_transaction_id_pk" PRIMARY KEY("transaction_id")
);

CREATE TABLE "entries" (
    "transaction_id" bigint NOT NULL,
	"account_id" int NOT NULL,
    "amount" numeric(7, 2) NOT NULL,
    
	CONSTRAINT "transactions_account_id_fk" FOREIGN KEY("account_id") REFERENCES accounts("account_id") ON DELETE CASCADE,
	CONSTRAINT "transactions_transaction_id_fk" FOREIGN KEY("transaction_id") REFERENCES transactions("transaction_id") ON DELETE CASCADE,
    CONSTRAINT "entries_amount_positive" CHECK (amount > 0)
);

-- +goose Down
DROP TABLE IF EXISTS "accounts" CASCADE;
DROP TABLE IF EXISTS "tags" CASCADE;
DROP TABLE IF EXISTS "account_tag" CASCADE;
DROP TABLE IF EXISTS "transactions" CASCADE;
DROP TABLE IF EXISTS "entries" CASCADE;
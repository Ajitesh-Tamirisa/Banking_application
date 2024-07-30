CREATE TABLE "accounts" (
  "account_number" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "transaction_id" bigserial PRIMARY KEY,
  "account_number" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "transfer_id" bigserial PRIMARY KEY,
  "from_account_number" bigint NOT NULL,
  "to_account_number" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("account_number");

CREATE INDEX ON "accounts" ("user_name");

CREATE INDEX ON "transactions" ("account_number");

CREATE INDEX ON "transfers" ("from_account_number");

CREATE INDEX ON "transfers" ("to_account_number");

ALTER TABLE "transactions" ADD FOREIGN KEY ("account_number") REFERENCES "accounts" ("account_number");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_number") REFERENCES "accounts" ("account_number");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_number") REFERENCES "accounts" ("account_number");

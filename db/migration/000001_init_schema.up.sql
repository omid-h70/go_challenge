CREATE TABLE "accounts" (
                            "account_id" bigserial PRIMARY KEY,
                            "owner" varchar NOT NULL,
                            "balance" bigint NOT NULL,
                            "currency" varchar NOT NULL,
                            "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "entries" (
                           "entry_id" bigserial PRIMARY KEY,
                           "account_id" bigserial,
                           "amount" bigint NOT NULL,
                           "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transfers" (
                             "transfer_id" bigint PRIMARY KEY,
                             "from_account_id" bigserial,
                             "to_account_id" bigserial,
                             "amount" bigint NOT NULL,
                             "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'it can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'it most be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("account_id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("account_id");
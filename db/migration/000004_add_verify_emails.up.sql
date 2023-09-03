CREATE TABLE "verify_emails" (
 "id" bigserial PRIMARY KEY,
 "username" varchar NOT NULL,
 "email" varchar NOT NULL,
 "secret_code" varchar NOT NULL,
 "is_used" bool NOT NULL DEFAULT (now()),
 "created_at" timestamptz NOT NULL DEFAULT (now()),
 "expired_at" timestamptz NOT NULL DEFAULT (now() + INTERVAL '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("user_name");

ALTER TABLE "users" ADD COLUMN "is_emailed_verified" bool NOT NULL DEFAULT false;
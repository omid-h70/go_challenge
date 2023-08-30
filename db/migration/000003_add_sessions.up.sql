CREATE TABLE sessions(
  "session_uuid" UUID PRIMARY KEY,
  "user_name" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT(now())
);

ALTER TABLE "sessions" ADD CONSTRAINT "fk_user_name" FOREIGN KEY("user_name") REFERENCES "users"("user_name");

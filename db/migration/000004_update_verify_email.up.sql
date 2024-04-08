

CREATE TABLE "verify_user_email" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
  "created_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "users" ADD is_email_verified bool NOT NULL DEFAULT false;
ALTER TABLE "verify_user_email" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

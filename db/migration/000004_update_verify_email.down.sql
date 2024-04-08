DROP TABLE IF EXISTS "verify_user_email" CASCADE;
ALTER TABLE "users" DROP COLUMN "is_email_verified";
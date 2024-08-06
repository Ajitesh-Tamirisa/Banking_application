ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "user_name_currency_key";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_user_name_fkey";

DROP TABLE IF EXISTS "users"
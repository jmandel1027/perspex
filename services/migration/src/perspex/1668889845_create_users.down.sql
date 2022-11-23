BEGIN;

DROP INDEX IF EXISTS users_email_uindex;

DROP INDEX IF EXISTS users_id_uindex;

DROP TABLE IF EXISTS "users";

END;

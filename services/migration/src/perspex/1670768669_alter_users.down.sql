
ALTER TABLE users DROP COLUMN IF EXISTS last_name;

ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name TEXT NOT NULL;

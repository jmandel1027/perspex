
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name TEXT NOT NULL;

ALTER TABLE users DROP COLUMN IF EXISTS full_name;
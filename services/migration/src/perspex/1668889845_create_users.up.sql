CREATE TABLE users (
  id         BIGSERIAL NOT NULL CONSTRAINT user_pk PRIMARY KEY,
  email      VARCHAR NOT NULL,
  full_name  VARCHAR NOT NULL,
  first_name VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_id_uindex
	ON "users" (id);

CREATE UNIQUE INDEX users_email_uindex
	ON "users" (email);

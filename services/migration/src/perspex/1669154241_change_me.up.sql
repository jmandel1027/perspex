CREATE TABLE IF NOT EXISTS organizations (
  id         BIGSERIAL NOT NULL CONSTRAINT organizations_pk PRIMARY KEY,
  name       VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS organizations_id_uindex
	ON "organizations" (id);

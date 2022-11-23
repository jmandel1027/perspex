-- migrate:down transaction:false

DROP INDEX CONCURRENTLY IF EXISTS organizations_id_uindex;

DROP TABLE IF EXISTS "organizations";

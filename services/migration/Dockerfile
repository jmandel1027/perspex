FROM migrate/migrate:v4.15.1 AS deployable

RUN apk add --no-cache bash

COPY src/ /mnt/src/

COPY bin/migration.sh /usr/local/bin/migration

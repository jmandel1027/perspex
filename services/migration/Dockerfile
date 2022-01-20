FROM migrate/migrate:v4.15.1 AS deployable

RUN apk add --no-cache bash

COPY db /mnt/db

COPY bin/entrypoint.sh /usr/local/bin/entrypoint
COPY bin/migration.sh /usr/local/bin/migration

ENTRYPOINT [ "entrypoint" ]
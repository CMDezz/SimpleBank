#!/bin/sh

set -e

echo "run db migration 2"
source /app/app.env
ls -l
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"

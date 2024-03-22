#!/bin/sh

set -e

echo "run db migration 2"
cat /app/app.env
source /app/app.env
ls -l ~/.app
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"

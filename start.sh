#!bin/bash

set -e

echo "run db migration"
# source env vars are replacing with docker vars
source /app/app.env
app/migrate -path /app/migration --database "$DB_SOURCE" -verbose up

echo "start the app"
echo "$@"
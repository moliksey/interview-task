#!/bin/bash
set -e
migrate -path ./schema -database 'postgres://postgres:root@interview-task-db:5432/postgres?sslmode=disable' up
>&2 echo "Postgres migrated"
exec "$@"
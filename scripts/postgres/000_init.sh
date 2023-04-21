#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER photoflux PASSWORD 'photoflux';
	CREATE DATABASE photoflux;
	GRANT ALL ON DATABASE photoflux TO photoflux;
EOSQL

psql_command() {
    psql --username photoflux --dbname photoflux -f "$1"
}

for f in /scripts/migrate/*.up.sql
do
    psql_command ${f}
done

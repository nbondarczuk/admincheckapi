#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER test;
	ALTER USER test WITH ENCRYPTED PASSWORD 'test';
	CREATE DATABASE argonadmindb;
	GRANT ALL PRIVILEGES ON DATABASE argonadmindb TO test;
EOSQL

#!/bin/bash
source .env

sleep 4 && goose -dir "./migrations" postgres "host=pg port=5432 dbname=$PG_DB_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable" up -v
sleep 4 && goose -dir "./migrations" postgres "host=pg-test port=5432 dbname=$PG_DB_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable" up -v

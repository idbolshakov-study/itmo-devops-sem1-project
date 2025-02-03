#!/bin/bash

DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="project-sem-1"
DB_USER="validator"
DB_PASSWORD="val1dat0r"

go mod download

PGPASSWORD=$DB_PASSWORD\
  psql \
    -h $DB_HOST -p $DB_PORT \
    -U $DB_USER -d $DB_NAME \
    -f sql/create_prices_table.sql

go build -o app main.go

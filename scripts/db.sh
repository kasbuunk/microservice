#!/usr/bin/env bash
set -e

source local.env

# specify function that creates database with given name

psql << EOF
DROP DATABASE ${SVC_DB_NAME}_test;
CREATE DATABASE ${SVC_DB_NAME}_test;

\c ${SVC_DB_NAME}_test;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
  id uuid DEFAULT gen_random_uuid() NOT NULL,
  email varchar(64) UNIQUE NOT NULL,
  password_hash varchar(64) UNIQUE NOT NULL,
  PRIMARY KEY(id)
);
EOF
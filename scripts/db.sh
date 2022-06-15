#!/usr/bin/env bash
set -e

source local.env

# specify function that creates database with given name

psql << EOF
DROP DATABASE ${SVC_DB_NAME};
DROP DATABASE ${SVC_DB_NAME}_test;
CREATE DATABASE ${SVC_DB_NAME};
CREATE DATABASE ${SVC_DB_NAME}_test;

\c ${SVC_DB_NAME};
CREATE TABLE users(
  id character(16),
  email character(64),
  password_hash character(64),
  PRIMARY KEY(id)
);

\c ${SVC_DB_NAME}_test;
CREATE TABLE users(
  id character(16),
  email character(64),
  password_hash character(64),
  PRIMARY KEY(id)
);


EOF
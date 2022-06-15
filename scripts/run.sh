#!/usr/bin/env bash
set -e

source local.env

make build

./bin/app

#!/bin/bash

set -e
docker rm -f postgres 2>/dev/null || true
make postgres
sleep 8
make createdb
make migrateup

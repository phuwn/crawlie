#!/bin/sh
sql-migrate up -config=dbconfig.yml -env=dev
crond -f
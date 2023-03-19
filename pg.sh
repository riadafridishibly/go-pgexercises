#!/bin/bash

PROGRAM="${0##*/}"

cmd() {
	echo "[#] $*" >&2
	"$@"
}

die() {
	echo "$PROGRAM: $*" >&2
	exit 1
}

PG_DOCKER_IMAGE=postgres:15.2-alpine3.17
PG_CONAINTER_NAME=pgexercies

docker_up() {
	cmd docker run --name "${PG_CONAINTER_NAME}" -p 5432:5432 -e POSTGRES_PASSWORD=password -d "${PG_DOCKER_IMAGE}"
	pg_populate_db
}

docker_down() {
	cmd docker stop "${PG_CONAINTER_NAME}"
	cmd docker rm "${PG_CONAINTER_NAME}"
}

download_sql_data() {
	cmd wget -O pgexercies.sql  https://pgexercises.com/dbfiles/clubdata.sql
}

pg_populate_db() {
	test -f ./pgexercies.sql || download_sql_data
	sleep 2 # wait for some time to initialize the db???
	cmd cat ./pgexercies.sql | docker exec -i "${PG_CONAINTER_NAME}" psql -U postgres
}

pg_query() {
	cmd docker exec -i "${PG_CONAINTER_NAME}" psql -U postgres -d exercises -c "${1}"
}

cmd_usage() {
	cat >&2 <<-_EOF
	Usage: $PROGRAM [ start | stop | query ]
	_EOF
}


if [[ $# -eq 1 && ( $1 == --help || $1 == -h || $1 == help ) ]]; then
	cmd_usage
elif [[ $# -eq 1 && $1 == start ]]; then
	docker_up
elif [[ $# -eq 1 && $1 == stop ]]; then
	docker_down
elif [[ $# -eq 2 && $1 == query ]]; then
	pg_query "$2"
else
	cmd_usage
	exit 1
fi

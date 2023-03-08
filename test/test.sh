#!/bin/bash

TEST_DIR=$(dirname "$0")

uwsgi --socket 127.0.0.1:3031 --wsgi-file "$TEST_DIR/foobar.py" &
UWSGI=$!

./caddy run --config "$TEST_DIR/Caddyfile" &
CADDY=$!

trap 'kill "$UWSGI" "$CADDY"' SIGINT SIGTERM EXIT

while ! nc -z localhost 3031 || ! nc -z localhost 9090; do
	true
done
curl -v http://127.0.0.1:9090/ -o "$TEST_DIR/response.txt" && grep -q Hello "$TEST_DIR/response.txt"

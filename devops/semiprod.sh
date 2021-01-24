#!/bin/bash

# Ensure we are in the right place.
curdir="$(basename "$(pwd)")"
if [ ! -f "project.conf" ]; then
	if [ "$curdir" == "devops" ]; then cd ..; curdir="$(basename "$(pwd)")"; fi
	if [ ! -f "project.conf" ]; then
		printf "Must be run out of the project directory.\n" >&2
		exit 1
	fi
fi

# Setup the on-exit cleanup
function kbye {
	if ! [ "$1" = "" ]; then
		echo "Killing $1"
		if ! kill $1 > /dev/null 2>&1; then
			sleep 2
			kill -9 $1 > /dev/null 2>&1
		fi
	fi
}

GOPID=""
SLPPID=""
function term {
	echo ""
	kbye "$GOPID"
	kbye "$SLPPID"

	echo "Cleanup finished."
}
trap term SIGHUP SIGINT SIGTERM

# Pull in configuration
source project.conf

# Make directories
if [ ! -d "app" ]; then
	mkdir app
fi

# Run an NPM build if needed
if [ -f "frontend/package.json" ]; then
	cd frontend
	npm i
	npm run build
	rm -r ../app/dist
	mv dist ../app
	cd ..
fi

# Start the go backend if this project has one
if [ -f "server/go.mod" ]; then
	cd server
	go build -o "../app/${curdir}.bin"
	cd ../app
	"./${curdir}.bin" &
	GOPID=$!
	cd ..
fi

# And wait for the user to send a ctrl-c or whatever.
sleep 5 &
SLPPID=$!
kill -STOP $SLPPID && wait $SLPPID

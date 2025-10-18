#!/bin/bash
cd "$(dirname "${BASH_SOURCE[0]}")"

if [ -z "$1" ]
  then
        echo "You must provide a binary name"
        exit 1
fi

cd got_bullet
echo "Got BUllet built binary is $1 and we are in $(eval pwd)"
#yea the binary needs to be build on the top level or whatever.
go build -buildvcs=false -o ../$1 ./cmd/got

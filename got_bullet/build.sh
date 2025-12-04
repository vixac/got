#!/bin/bash
cd "$(dirname "${BASH_SOURCE[0]}")"

go build -buildvcs=false -o got

#!/bin/bash
cd "$(dirname "${BASH_SOURCE[0]}")"
CMD="go test -v ./..."
if [ -z "$(which grc)" ]
then 
   echo "No grc install. No color"
   eval $CMD
else 
   grc  $CMD
fi

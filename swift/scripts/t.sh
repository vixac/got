#!/bin/bash
if [[ -z "$GOT" ]]; then
   GOT=~/.got
fi

P=$GOT/contents/free
SCRIPT=$GOT/scripts/open_then_timestamp.sh

if [ $# -eq 0 ]; then
  $SCRIPT $P/t.txt
  exit
fi

  
if [ $# -eq 1 ]; then
  $SCRIPT $P/$1.txt
  exit
fi

if [ $# -eq 2 ]; then
  mkdir -p $P/$1
  $SCRIPT $P/$1/$2.txt
  exit
fi


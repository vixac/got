if [[ -z "$GOT" ]]; then
   GOT=~/.got
fi
. $GOT/got_env
grep -l $1 $GOT_ACTIVE/*_summary.got > $GOT_OUTPUT_FILE

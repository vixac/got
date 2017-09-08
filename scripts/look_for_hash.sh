if [[ -z "$GOT" ]]; then
   GOT=~/.got
fi
grep -l $1 $GOT/contents/active/*_summary.got > $GOT/.tmpdata

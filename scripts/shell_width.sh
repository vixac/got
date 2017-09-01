if [[ -z "$GOT" ]]; then
   GOT=~/.got
fi
. $GOT_SRC/got_env
tput cols > $GOT_OUTPUT_FILE

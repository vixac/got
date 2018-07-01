#!/bin/bash
#opens $1 filename in vim, then when closed, inserts the timestamp at the top  of the file.
FILE=$1
TMPFILE=~/.got_tmp.txt
touch $FILE
#copying file to tmp so that we can see if it changes
cp $FILE $TMPFILE

#take the time editing began.
DATE=$('date')
vim -c 'startinsert' $FILE

#take a diff to see if we've made a change.
DIFF=$(diff $FILE $TMPFILE)
rm $TMPFILE
if [ "$DIFF" != "" ] 
then
    # echo newline with the date, then cat it with the rest of the file to 
    LINE="----------------------------"
    (echo -e '\n' && echo "$LINE" && echo $DATE && echo -e "$LINE \n") | cat - $FILE > $TMPFILE && mv $TMPFILE $FILE
fi

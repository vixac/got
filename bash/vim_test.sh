#!/bin/bash

FILE=$1.txt
TMPFILE=.$FILE.vxdaytmp
touch $FILE
#copying file to tmp so that we can see if it changes
cp $FILE $TMPFILE

#take the time editing began.
DATE=$('date')
vim $FILE

#take a diff to see if we've made a change.
DIFF=$(diff $FILE $TMPFILE)
rm $TMPFILE
if [ "$DIFF" != "" ] 
then
    # echo newline with the date, then cat it with the rest of the file to 
    (echo -e '\n' && echo $DATE) | cat - $FILE > $TMPFILE && mv $TMPFILE $FILE
fi

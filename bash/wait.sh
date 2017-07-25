echo "starting timer for list $1 and hash $2"
FILE=$1
HASH=$2

START_TIME=$(date +"%Y-%m-%dT%H:%M:%S")

function  finish {
  END_TIME=$(date +"%Y-%m-%dT%H:%M:%S")
  #echo "finished.TODO send $START_TIME AND $END_TIME to $FILE under hash $HASH"
  TOKEN_LINE=$(echo "-> $HASH $START_TIME $END_TIME")
  #echo line is $TOKEN_LINE
  #echo writing to file $FILE
  echo $TOKEN_LINE >> $FILE
  #$($VXDAY2_SRC_DIR/bash/append.sh $TOKEN_LINE $FILE)
  exit
}
trap finish EXIT
while read line
do
  echo "LINE IS $line, TODO THIS DOESNT WORK."
done 



echo "WAIT WAS CALLED...."
FILE=$1
HASH=$2

START_TIME=$(date +"%Y-%m-%dT%H:%M:%S")

function  finish {
  END_TIME=$(date +"%Y-%m-%dT%H:%M:%S")
  TOKEN_LINE=$(echo "-> $HASH $START_TIME $END_TIME")
  echo $TOKEN_LINE >> $FILE
  exit
}
trap finish EXIT
while read line
do
  echo ""
done 



SWIFTC_EXISTS=$(which swiftc)
if [[ -z $SWIFTC_EXISTS ]]; then 
   echo "Error, you don't have the swift compiler installed yet. You'll need to download it"
   exit 1
fi

BASE=./swift
echo "compiling got.."
swiftc $BASE/main.swift $BASE/VxdayTable.swift $BASE/Trap.swift $BASE/VxdayTypes.swift $BASE/VxdayRead.swift $BASE/VxdayView.swift $BASE/VxdayInstruction.swift $BASE/VxdayExec.swift $BASE/VxdayUtil.swift $BASE/VxdayInput.swift  -o got
echo "compilation complete!"

if [[ -z "$GOT" ]]; then
   GOT=~/.got
fi


mkdir -p $GOT/contents/active
mkdir -p $GOT/contents/retired
mkdir -p $GOT/scripts
cp scripts/*.sh $GOT/scripts
chmod +x $GOT/scripts/*.sh

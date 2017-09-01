
echo "GOT SOURCE IS '$GOT_SRC'"
if [[ -z "$GOT_SRC" ]]; then
 echo "Error, GOT_SRC is not set. It's the root dir of the Got git repo."
 exit 1
fi
SWIFTC_EXISTS=$(which swiftc)
if [[ -z $SWIFTC_EXISTS ]]; then 
   echo "Error, you don't have the swift compiler installed yet. You'll need to download it"
   exit 1
fi

BASE=$GOT_SRC/swift
echo "compiling got.."
swiftc $BASE/main.swift $BASE/VxdayTable.swift $BASE/Trap.swift $BASE/VxdayTypes.swift $BASE/VxdayRead.swift $BASE/VxdayView.swift $BASE/VxdayInstruction.swift $BASE/VxdayExec.swift $BASE/VxdayUtil.swift $BASE/VxdayInput.swift  -o $GOT_SRC/got
echo "compilation complete!"
echo "moving got to /usr/local/bin"
mv $GOT_SRC/got /usr/local/bin

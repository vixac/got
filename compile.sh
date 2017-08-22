if [[ -z "${GOT_SRC}" ]]; then
 echo "Error, GOT_SRC is not set. It is the root dir of the Got git repo."
fi
BASE=$GOT_SRC/swift
swiftc $BASE/main.swift $BASE/VxdayTable.swift $BASE/Trap.swift $BASE/VxdayTypes.swift $BASE/VxdayRead.swift $BASE/VxdayView.swift $BASE/VxdayInstruction.swift $BASE/VxdayExec.swift $BASE/VxdayUtil.swift $BASE/VxdayInput.swift  -o $GOT_SRC/got


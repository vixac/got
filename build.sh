BASE=./swift
echo "compiling got.."
swiftc $BASE/main.swift $BASE/VxdayTable.swift $BASE/Trap.swift $BASE/VxdayTypes.swift $BASE/VxdayRead.swift $BASE/VxdayView.swift $BASE/VxdayInstruction.swift $BASE/VxdayExec.swift $BASE/VxdayUtil.swift $BASE/VxdayInput.swift  -o got
echo "compilation complete!"


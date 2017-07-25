
#cat VxdayInstruction.swift VxdayExec.swift VxdayUtil.swift VxdayParser.swift test.swift  | swift - $@
cat VxdayInstruction.swift VxdayExec.swift VxdayUtil.swift VxdayParser.swift test.swift  > flattened.swift

swift flattened.swift $0

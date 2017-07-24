//
//  VxdayExec.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class VxdayExec {
    
    
    private static let activeDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY_ACTIVE_DIR")!
    }()
    
    private static let retiredDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY_RETIRED_DIR")!
    }()
    
    private static let bashDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY_SRC_DIR")! + "/bash"
    }()
    private static let starVxday = "_*.vxday"
    
    @discardableResult
    static func shell(_ args: String...) -> Int32 {
        
        print("about to do this: \(args)")
        let task = Process()
        task.launchPath = "/usr/bin/env"
        task.arguments = args
        task.launch()
        task.waitUntilExit()
        return task.terminationStatus
    }
    
    static func getEnvironmentVar(_ name: String) -> String? {
        guard let rawValue = getenv(name) else { return nil }
        return String(utf8String: rawValue)
    }

    
    static func retire(_ list: ListName) {
        let script = VxdayExec.bashDir + "/retire.sh"
        VxdayExec.shell(script, list.name)
    }
    
    static func unretire(_ list: ListName) {
        let script = VxdayExec.bashDir + "/unretire.sh"
        VxdayExec.shell(script, list.name)
    }
    
}

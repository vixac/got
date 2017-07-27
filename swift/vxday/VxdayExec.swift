//
//  VxdayExec.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

enum FileType : String {
    case summary = "summary"
    case tokens  = "tokens"
    case note = "note"
}

enum Script : String {
    case retire = "retire.sh"
    case unretire = "unretire.sh"
    case append = "append.sh"
    case removeLine = "remove_line.sh"
    case note = "note.sh"
    case wait = "wait.sh"
}

class VxdayFile {
    
    static let bashDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY2_SRC_DIR")! + "/bash"
    }()
    
    static let activeDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY2_ACTIVE_DIR")!
    }()
    
    static let retiredDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY2_RETIRED_DIR")!
    }()
    
    static func getScriptPath(_ script: Script) -> String {
        return VxdayFile.bashDir + "/" + script.rawValue
    }
    
    static func getSummaryFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_summary.vxday"
    }
    
    static func getNoteFilename(_ list: ListName, hash: Hash) -> String {
        let end = "_" + hash.hash + ".vxday"
        return VxdayFile.activeDir + "/" + list.name + end
    }
    
    static func getTokenFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_tokens.vxday"
    }
 
}

class VxdayExec {
    
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
 
    static func lessList(_ list: ListName) {
        let files = VxdayFile.activeDir + "/" + list.name + "_*.vxday"
        VxdayExec.shell("cat", files)
    }
    static func what() {
        let lists = VxdayReader.allLists()
        var allItemsEver: [Item] = []
        lists.forEach { list in
            allItemsEver += VxdayReader.itemsInList(list)
        }
        let view = VxdayView(allItemsEver)
        view.oneLiners().forEach {print($0)}
    }
    
    static func x(_ hash: Hash) {
        print("TODO X THIS HASH: \(hash)")
    }
    
    static func all() {
        //all lists and their summaries.
        let lists = VxdayReader.allLists()
        var allItemsEver: [Item] = []
        lists.forEach { list in
            allItemsEver += VxdayReader.itemsInList(list)
        }
        let view = VxdayView(allItemsEver)
        view.renderAll().forEach { print($0) }
    }
    
    static func allList(_ list: ListName) {
        let items = VxdayReader.itemsInList(list)
        let view  = VxdayView(items)
        let jobsStrings = view.renderAll()
        jobsStrings.forEach { print($0)}
    }
    
    //TODO try to write these using mv
    static func retire(_ list: ListName) {
        let script = VxdayFile.getScriptPath(.retire)
        VxdayExec.shell(script, list.name)
    }
    
    //TODO try to write these using mv
    static func unretire(_ list: ListName) {
        let script = VxdayFile.getScriptPath(.unretire)
        VxdayExec.shell(script, list.name)
    }
    
    static func append(_ list: ListName, content: String ) {
        let script = VxdayFile.getScriptPath(.append)
        let filename = VxdayFile.getSummaryFilename(list)
        print("about to run \(script) on \(filename) with content: \(content)")
        VxdayExec.shell(script, content, filename)
    }
    
    static func note(_ list: ListName, hash: Hash) {
        let script = VxdayFile.getScriptPath(.note)
        let filename = VxdayFile.getNoteFilename(list, hash: hash)
        VxdayExec.shell(script, filename)
        
    }
    
    static func wait(_ list: ListName, hash: Hash) {
        let script = VxdayFile.getScriptPath(.wait)
        let filename = VxdayFile.getTokenFilename(list)
        VxdayExec.shell(script, filename, hash.hash)
    }
    
}

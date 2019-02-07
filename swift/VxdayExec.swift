//
//  VxdayExec.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

enum Script : String {
    case retire = "retire.sh"
    case unretire = "unretire.sh"
    case append = "append.sh"
    case note = "note.sh"
    case wait = "wait.sh"
    case lookForHash = "look_for_hash.sh"
    case removeLine = "remove_line.sh"
    case shellWidth = "shell_width.sh"
    case vim = "open_then_timestamp.sh"
    case exec = "exec.sh"
}

class VxdayFile {
    
    static let gotBase: String = {
        if let s =  VxdayExec.getEnvironmentVar("GOT") {
             return s 
}
       else { 
      let home = VxdayExec.getEnvironmentVar("HOME") ?? "~" 
      return home + "/.got"
}
    }()
    static let contentDir: String = {
       return   VxdayFile.gotBase + "/contents"
    }()
    static let bashDir: String = {
       return  VxdayFile.gotBase + "/scripts"
    }()
    
    static let activeDir: String = {
       return  VxdayFile.contentDir + "/active"
    }()
    
    static let retiredDir: String = {
       return  VxdayFile.contentDir + "/retired"
    }()
    
    static let outputFile : String = {
       return  VxdayFile.gotBase + "/.tmpdata"
    }()
    
    static func getScriptPath(_ script: Script) -> String {
       return VxdayFile.bashDir + "/" + script.rawValue
    }
    
    static func getSummaryFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_summary.got"
    }
    static func getNoteFilename(_ list: ListName, hash: Hash) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_" + hash.hash + "_notes.got"
    }
    
    static func getCompleteFilename(_ list: ListName) -> String  {
        return VxdayFile.activeDir + "/" + list.name + "_complete.got"
    }
    static func getTokenFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_tokens.got"
    }

}


struct VxdayExecSessionGlobals {
    let start: CreationDate
    let hash: Hash
    let list : ListName
}
class VxdayExec {
    
    static var globals : VxdayExecSessionGlobals? = nil
    private static let starVxday = "_*.got"
    
    @discardableResult
    static func shell(_ args: String...) -> Int32 {
        let task = Process()
        task.launchPath = "/usr/bin/env"
        task.arguments = args
        task.launch()
        task.waitUntilExit()
        return task.terminationStatus
    }
    /*
    static func run(_ pName: String, _ items: [String]) {
        print("VX: attempting to launch: \(pName) with itemss: \(items)")
        let process = Process()
        if #available(OSX 10.13, *) {
            process.executableURL = URL(fileURLWithPath:pName)
            process.arguments = ["-la"]
            
            process.terminationHandler = { (process) in
                print("\ndidFinish: \(!process.isRunning)")
            }
            do {
                try process.run()
            } catch let e  {
                print("VX: is this workign? \(e)")
            }
        } else {
            print("VX: running processes is not available before Mac OS 10.13")
            // Fallback on earlier versions
        }
        
    }
 */
    
    static func quickRun(_ name: String, args: [String]) {
        if #available(OSX 10.13, *) {
       // let url = URL(fileURLWithPath:name)
        do {
            let process = Process()
           // let pipe = Pipe()
            process.executableURL = URL(fileURLWithPath:name)
            //process.standardOutput = pipe
            //process.standardError = pipe
            process.arguments = args
           // process.waitUntilExit()
            process.terminationHandler =  { (process) in
                print("\ndidFinish: \(!process.isRunning)")
            }
            try process.run()
        } catch let e  {
            print("VX: is this workign? \(e)")
            }
        } else {
             print("VX: running processes is not available before Mac OS 10.13")
        }
    }
     @discardableResult static func shellNew(_ launchPath: String, _ arguments: [String] = []) -> (String? , Int32) {
        let task = Process()
        
        task.launchPath = launchPath
        task.arguments = arguments
        
        let pipe = Pipe()
        
        task.standardOutput = pipe
        task.standardError = pipe
        task.launch()
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        let output = String(data: data, encoding: .utf8)
        task.waitUntilExit()
        
        return (output, task.terminationStatus)
    }
    
    static func getEnvironmentVar(_ name: String) -> String? {
        guard let rawValue = getenv(name) else { 
           return nil
        }
        return String(utf8String: rawValue)
    }
 
    
    static func openTop(_ list: ListName?) {
        print("VX: opentop called")
        if let l = list  {
            print("VX: opening  shell new top for: \(l)")
           // let echo = "/bin/echo"
            let vim = "/usr/bin/vim"
            //let vimScript = VxdayFile.getScriptPath(.vim)
            let execScript = VxdayFile.getScriptPath(.exec)
            //VxdayExec.quickRun("/usr/bin/vim", args:  ["wtf.txt"]) //this opens then closes instantly
            //VxdayExec.quickRun(execScript, args:  ["vim test.txt"]) //this opens then closes instantly
            VxdayExec.shell(execScript, vim)
        }
    }
    
    static func lessList(_ list: ListName) {
        let files = VxdayFile.activeDir + "/" + list.name + "_*.got"
        VxdayExec.shell("cat", files)
    }
    
    static func hashToListName(_ hash: Hash) -> ListName? {
        guard hash.isValid() else {
            print("This isn't a valid hash: \(hash)")
            return nil
        }
        let findHashScript = VxdayFile.getScriptPath(.lookForHash)
        VxdayExec.shell(findHashScript, hash.hash)
        guard let listFile = getOutput() else {
            print("This hash is not active: \(hash.hash)")
            return nil
        }
        guard let listName =  VxdayUtil.beforeUnderscore(VxdayUtil.afterHyphen(listFile)) else {
            print("Error reading the list summary file: \(listFile)")
            return nil
        }
        return ListName(listName)
    }
    
    static func getOutput() -> String? {
        return VxdayReader.readFile(VxdayFile.outputFile).first
    }
    
    static func getShellWidth() -> Int {
        VxdayExec.shell(VxdayFile.getScriptPath(.shellWidth))
        guard let str = getOutput() else {
            print("Error getting shell width")
            return 150
        }
        guard let width = Int(str) else {
            print("Error, width is not a number")
            return 150
        }
        return width - 30 // need the 30 and don't really know why
    }
    
    static func removeActiveHash(_ hash: Hash, list: ListName?) {
        guard let l = getListIfNeeded(hash, list: list) else {
            print("Error: Can't find list for hash: \(hash.hash)")
            return
        }
        
        let summaryFileName = VxdayFile.getSummaryFilename(l)
        let scriptFile = VxdayFile.getScriptPath(.removeLine)
        
        VxdayExec.shell(scriptFile, hash.hash, summaryFileName)
    }
    
    static func what() {
        let lists = VxdayReader.allLists()
        var allItemsEver: [Item] = []
        lists.forEach { list in
            allItemsEver += VxdayReader.itemsInList(list)
        }
        WhatView(allItemsEver).toTable(getShellWidth()).render().forEach {print($0)}
    }
    
    
    static func getListIfNeeded(_ hash: Hash,  list: ListName?) -> ListName? {
        if let _ = list {
            return list
        }
        return hashToListName(hash)
    }
    
    static func showNotes(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            return
        }
        let filename = VxdayFile.getNoteFilename(list, hash: hash)
        VxdayExec.shell("cat", filename)
        
    }
    
    static func hashToJobOrTask(_ hash: Hash, list: ListName?) ->  Item? {
        var theList = list
        if theList == nil {
            theList  = hashToListName(hash)
        }
        guard let l = theList else {
            print("Error finding list for hash: \(hash.hash)")
            return nil
        }
        
        // that job or task thing is just in case we have token or notes in the same file
        let items = VxdayReader.itemsInList(l).filter { $0.vxItem().hash.hash == hash.hash && ( $0.vxItem().itemType() == .job || $0.vxItem().itemType() == .task) || $0.vxItem().itemType() == .now }
        guard let item  = items.first  else {
            print("Dev Error finding hash \(hash.hash) in list \(l.name)")
            return nil
        }
        return item
    }
    
    static func info(_ hash: Hash) {
        print("TDOO info on \(hash)")
        let lineCount = noteLineCount(hash)
        print("GOT LINE COUNT: \(lineCount)")
    }
    static func note(_ hash: Hash) {
        /*
        guard let list = hashToListName(hash) else {
            print("Error finding list name for hash: \(hash)")
            return
        }
 */
        print("TODO this doesnt woork yet.")
        
    }
    static func showComplete(_ list: ListName?) {
        if let l = list {
            let completeFile = VxdayFile.getCompleteFilename(l)
            let contents = VxdayReader.readFile(completeFile)
            let items = VxdayReader.linesToItems(contents, list: l)
            let sortedItems = items.sorted(by: {$0.vxItem().creation.date < $1.vxItem().creation.date})
            
            
            let view = CompleteTableView(sortedItems)
            let table = view.toTable(getShellWidth())
            table.render().forEach {
                print($0)}
        
        }
        else {
            var allCompleteItems: [Item] = []
            
            let all = VxdayReader.allLists()
            all.forEach { list in
                allCompleteItems += VxdayReader.completeItemsInList(list)
            }
            
            let view = CompleteTableView(allCompleteItems)
            let table = view.toTable(getShellWidth())
            table.render().forEach {print($0)}
        }
    }
    
    static func waitForUser(cb:( (CreationDate,CompletionDate) -> Void), note: (String) -> Void) {
        let start  = VxdayUtil.now()
        Trap.handle(signal: Trap.Signal.interrupt) { signal in
            VxdayExec.bypassCPointerContextIssueByUsingStaticStateToSaveToken()
        }
        while let line = readLine(strippingNewline: true) {
            if line == "stop" {
                let end = VxdayUtil.now()
                cb(CreationDate(start), CompletionDate(end))
                return
            }
            else {
                
                note(line)
            }
        }
    }
    static func noteLineCount(_ hash: Hash) -> Int {
        guard let list = hashToListName(hash) else {
            print("No hash found: \(hash)")
            return 0
        }
        let filename  = VxdayFile.getNoteFilename(list, hash: hash)
        let count = VxdayReader.readFile(filename).count
        return count
    }
    
    static func takeNote(_ list: ListName, hash: Hash, note: String) {
        let noteFile = VxdayFile.getNoteFilename(list, hash: hash)
        let script = VxdayFile.getScriptPath(.append)
        VxdayExec.shell(script, note, noteFile)
    }
    
    static func bypassCPointerContextIssueByUsingStaticStateToSaveToken() {
        guard let g = VxdayExec.globals else {
            print("Dev error: It appears an interuption doesnt have the static globals in place")
            return
        }
        VxdayExec.saveToken(g.list, hash: g.hash, creation: g.start, completion: CompletionDate(VxdayUtil.now()))
    }
    
    static func saveToken(_ list: ListName, hash: Hash, creation: CreationDate, completion: CompletionDate) {
        let token = VxToken(list: list, hash: hash, creation: creation, completion: completion)
        self.storeItem(token)
        let timeBreakdown = TimeBreakdown(start: creation.date, end: completion.date)
        OneLinerView.showTimerStopped(timeBreakdown).render().forEach{ print($0)}
        
    }
    
    static func report(_ days: IntOffset, list: ListName?) {
        
        let allLists = list == nil ? VxdayReader.allLists() : [list!]
        
        let date = TokenDayView.dateOfStart(daysAgo: days)
        let report = TokenDayView()
        VxdayReader.allTokensAfter(date, lists: allLists).forEach {
            //TODO dont actually need to do this for each token, because its per hash per day that we need it but oh well.
            let description = VxdayReader.descriptionForHash($0.hash, list: $0.list)
            var t = $0
            t.description = description
            report.addToken(t)
        }
     
        report.toTable(days, width: getShellWidth()).render().forEach { print($0)}
    }
    
    static func startTokenSession(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            print("Error finding this hash in an active list: \(hash.hash)")
            return
        }
        let now = VxdayUtil.now()
        VxdayExec.takeNote(list, hash: hash, note: VxdayUtil.datetimeFormatter.string(from: now))
        VxdayExec.globals = VxdayExecSessionGlobals(start: CreationDate(now), hash: hash, list: list)
        let description = VxdayReader.descriptionForHash(hash , list: list)
        OneLinerView.showTimerStarted(list, time: now, hash: hash, description: description ?? Description("")).render().forEach { print($0)}
        waitForUser (cb: { start, end in
            saveToken(list, hash: hash , creation: start , completion: end)
        }, note: { text in
            VxdayExec.takeNote(list, hash: hash, note: text)
        })
        
    }
    static func remove(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            print("erorr finding hash: \(hash)")
            return
        }
        guard let description = VxdayReader.descriptionForHash(hash, list: list) else {
            print("Error getting description for list: '\(list) and hash \(hash)")
            return
        }
        OneLinerView.showHashRemoved(hash, list: list, description: description).render().forEach {print($0)}
        removeActiveHash(hash, list: list)
    }
    
    static func createTask(_ list: ListName, description: Description ) {
        if !VxdayReader.isListPresent(list) {
            OneLinerView.showNewList(list).render().forEach { print($0)}
        }
        
        let now = VxdayUtil.now()
        let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
        let created = CreationDate(now)
        
        let vxtask = VxTask(list: list, hash: hash, creation: created, description: description, completion: nil)
        VxdayExec.storeItem(vxtask)

        OneLinerView.showItemCreatedOneLiner(vxtask).render().forEach {print($0)}
    }
    
    static func createNow(_ list: ListName, description: Description) {
        if !VxdayReader.isListPresent(list) {
            OneLinerView.showNewList(list).render().forEach { print($0)}
        }
        let now = VxdayUtil.now()
        let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
        let created = CreationDate(now)
        
        let nowItem = VxNow(creation: created, list: list, hash: hash, completion: nil, description: description)
        VxdayExec.storeItem(nowItem)
        OneLinerView.showItemCreatedOneLiner(nowItem).render().forEach {print($0)}
    }
    
    static func createJob(_ list: ListName, offset: IntOffset, description: Description ) {
        
        if !VxdayReader.isListPresent(list) {
            OneLinerView.showNewList(list).render().forEach { print($0)}
        }
        let now = VxdayUtil.now()
        let created = CreationDate(now)
        let deadline = DeadlineDate(VxdayUtil.increment(date: now, byDays: offset.offset))
        let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
        let vxjob = VxJob(list: list, hash: hash, creation: created, deadline: deadline, description: description , completion: nil)
        
        VxdayExec.storeItem(vxjob)
        OneLinerView.showItemCreatedOneLiner(vxjob).render().forEach {print($0)}
    }
    
 
    static func x(_ hash: Hash) {
        
        guard let list = hashToListName(hash) else {
            return
        }
        guard let item = hashToJobOrTask(hash, list: list) else {
            print("Error extracting job from hash: \(hash)")
            return
        }
        
        removeActiveHash(hash, list: list)
        let completeItem = item.vxItem().complete()
        storeItem(completeItem)
        
        var description: Description? = nil
        if let j = item.getJob() {
            description = j.description
        }
        if let t = item.getTask() {
            description = t.description
        }
        if let n = item.getNow() {
            description = n.description
        }
        guard let d = description else {
            print("No description found.")
            return
        }
        OneLinerView.showHashCompleted(hash, list: list, description: d).render().forEach{print($0)}
 
    }
    
    static func all() {
        
        let lists = VxdayReader.allLists()
        var allItemsEver: [Item] = []
        lists.forEach { list in
            allItemsEver += VxdayReader.itemsInList(list)
        }
        
        let view  = ItemTableView(allItemsEver)
        let width = getShellWidth()
        let table = view.toTable(width)
        table.render().forEach { print($0) }
    }
    
    static func allList(_ prefix: String) {
        let allLists = VxdayReader.allLists(prefix)
        var allItems: [Item] = []
        allLists.forEach {  list in
                allItems += VxdayReader.itemsInList(list)
        }
        
        let view  = ItemTableView(allItems)
        let table = view.toTable(getShellWidth())
        table.render().forEach { print($0) }
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
    
    
    static func storeItem(_ item: VxItem) {
        let script = VxdayFile.getScriptPath(.append)
        let list  = item.list
        var filename : String = ""
        if item.itemType() == ItemType.token {
            filename = VxdayFile.getTokenFilename(list)
        }
        else if item.isComplete() {
            filename = VxdayFile.getCompleteFilename(list)
        }
        else {
            filename = VxdayFile.getSummaryFilename(list)
        }
        let content = item.toVxday()
        VxdayExec.shell(script, content, filename)
        
    }
    
    static func gotInfo() {
        let gotDir = VxdayFile.gotBase
        print("Your got base is: \(gotDir)")
    }
    
    static func help() {
        HelpView.toTable(getShellWidth()).render().forEach {print($0)}
    }
    
    static func wait(_ list: ListName, hash: Hash) {
        let script = VxdayFile.getScriptPath(.wait)
        let filename = VxdayFile.getTokenFilename(list)
        VxdayExec.shell(script, filename, hash.hash)
    }
    
}

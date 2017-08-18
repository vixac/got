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
}

class VxdayFile {
    
    static let bashDir: String = {
        VxdayExec.getEnvironmentVar("GOT_SRC")! + "/bash"
    }()
    
    static let activeDir: String = {
        VxdayExec.getEnvironmentVar("GOT_ACTIVE")!
    }()
    
    static let retiredDir: String = {
        VxdayExec.getEnvironmentVar("GOT_RETIRED")!
    }()
    
    static let outputFile : String = {
        VxdayExec.getEnvironmentVar("GOT_OUTPUT_FILE")!
    }()
    
    static func getScriptPath(_ script: Script) -> String {
        return VxdayFile.bashDir + "/" + script.rawValue
    }
    
    static func getSummaryFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_summary.vxday"
    }
    static func getNoteFilename(_ list: ListName, hash: Hash) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_" + hash.hash + "_notes.vxday"
    }
    
    static func getCompleteFilename(_ list: ListName) -> String  {
        return VxdayFile.activeDir + "/" + list.name + "_complete.vxday"
    }
    static func getTokenFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_tokens.vxday"
    }

}


struct VxdayExecSessionGlobals {
    let start: CreationDate
    let hash: Hash
    let list : ListName
}
class VxdayExec {
    
    static var globals : VxdayExecSessionGlobals? = nil
    private static let starVxday = "_*.vxday"
    
    @discardableResult
    static func shell(_ args: String...) -> Int32 {
        let task = Process()
        task.launchPath = "/usr/bin/env"
        task.arguments = args
        task.launch()
        task.waitUntilExit()
        return task.terminationStatus
    }
    
    
    @discardableResult
    static func shellNoWait(_ args: String...) -> Int32 {
        let task = Process()
        task.launchPath = "/usr/bin/env"
        task.arguments = args
        task.launch()
        return task.terminationStatus
    }
    
    static func getEnvironmentVar(_ name: String) -> String? {
        guard let rawValue = getenv(name) else { 
           print("Error, can't find environment variable: \(name). You've probably not set GOT_SRC or GOT_BASE, or not sourced the got_env file")
           return nil
        }
        return String(utf8String: rawValue)
    }
 
    static func lessList(_ list: ListName) {
        let files = VxdayFile.activeDir + "/" + list.name + "_*.vxday"
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

        let view = ItemView(allItemsEver)
        let buckets = ItemView.toBuckets(view.items)
        
        view.oneLiners(buckets).forEach {print($0)}
        
        let global = view.globalOneLiner(buckets: buckets)
        print("-----------------------------------------------------------------------")
        print((global) )
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
        let items = VxdayReader.itemsInList(l).filter { $0.vxItem().hash.hash == hash.hash && ( $0.vxItem().itemType() == .job || $0.vxItem().itemType() == .task) }
        guard let item  = items.first  else {
            print("Dev Error finding hash \(hash.hash) in list \(l.name)")
            return nil
        }
        return item
    }
    
    static func note(_ hash: Hash) {
        /*
        guard let list = hashToListName(hash) else {
            print("Error finding list name for hash: \(hash)")
            return
        }
 */
        print("TODO this doesnt woork yet.")
        //let script = VxdayFile.getScriptPath(.note)
        //let filename = VxdayFile.getNoteFilename(list, hash: hash)
       // VxdayExec.shellNoWait ("vim", filename)doesnt work.
        
    }
    static func showComplete(_ list: ListName?) {
        if let l = list {
            let completeFile = VxdayFile.getCompleteFilename(l)
            let contents = VxdayReader.readFile(completeFile)
            let items = VxdayReader.linesToItems(contents, list: l)
            
            
            let view = CompleteTableView(items)
            let table = view.toTable()
            table.render().forEach {
                print("rendering $0 which is \($0)")
                print($0)}
        
        }
        else {
            var allCompleteItems: [Item] = []
            
            let all = VxdayReader.allLists()
            all.forEach { list in
                allCompleteItems += VxdayReader.completeItemsInList(list)
            }
            
            let view = CompleteTableView(allCompleteItems)
            let table = view.toTable()
            table.render().forEach {print($0)}
        }
    }
    
    static func waitForUser(cb:( (CreationDate,CompletionDate) -> Void), note: (String) -> Void) {
        let start  = VxdayUtil.now()
        Trap.handle(signal: Trap.Signal.interrupt) { signal in
            VxdayExec.bypassCPointerContextIssueByUsingStaticStateToSaveToken()
        }
        while let line = readLine(strippingNewline: true) {
            if line == "done" {
                let end = VxdayUtil.now()
                cb(CreationDate(start), CompletionDate(end))
                return
            }
            else {
                
                note(line)
            }
        }
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
        let view = ItemView(Item.token(token))
        view.renderAll().forEach { print($0)}
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
     
        report.toTable(days).render().forEach { print($0)}
    }
    
    static func startTokenSession(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            print("Error finding this hash in an active list: \(hash.hash)")
            return
        }
        guard let item = hashToJobOrTask(hash, list: list) else {
            print("Error finding item for hash \(hash) in list \(list)")
            return
        }
        let now = VxdayUtil.now()
        VxdayExec.globals = VxdayExecSessionGlobals(start: CreationDate(now), hash: hash, list: list)
        
        let description = item.getJob() != nil  ? item.getJob()!.description : item.getTask()!.description
        OneLinerView.showTimerStarted(list, time: now, hash: hash, description:  description).render().forEach { print($0)}
        waitForUser (cb: { start, end in
            saveToken(list, hash: hash , creation: start , completion: end)
        }, note: { text in
            let now = VxdayUtil.datetimeFormatter.string(from: VxdayUtil.now())
            VxdayExec.takeNote(list, hash: hash, note: "\(now): \(text)")
        })
        
    }
    static func remove(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            return
        }
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
        guard let item = hashToJobOrTask(hash, list: list)?.vxItem() else {
            print("Error extracting job from hash: \(hash)")
            return
        }
        removeActiveHash(hash, list: list)
        let completeItem = item.complete()
        storeItem(completeItem)
        
    }
    
    static func all() {
        
        let lists = VxdayReader.allLists()
        var allItemsEver: [Item] = []
        lists.forEach { list in
            allItemsEver += VxdayReader.itemsInList(list)
        }
        
        let view  = ItemTableView(allItemsEver)
        let table = view.toTable()
        table.render().forEach { print($0) }
    }
    
    static func allList(_ list: ListName) {
        
        let items = VxdayReader.itemsInList(list)
        
        let view  = ItemTableView(items)
        let table = view.toTable()
        table.render().forEach { print($0) }
        //let view  = ItemView(items)
        //let jobsStrings = view.renderAll()
        //jobsStrings.forEach { print($0)}
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
    
    
    
    static func wait(_ list: ListName, hash: Hash) {
        let script = VxdayFile.getScriptPath(.wait)
        let filename = VxdayFile.getTokenFilename(list)
        VxdayExec.shell(script, filename, hash.hash)
    }
    
}

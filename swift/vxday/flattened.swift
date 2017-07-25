//
//  VxdayInstruction.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

//TODO make this config
enum ItemType : String {
    case complete = "x."
    case tokenEntry = "->."
    case job = "=."
    case task = "-."
}


protocol LineItem {
    var hash: Hash {get}
    var description: Description {get}
    func toString(complete: Bool) -> String
    func itemType() -> ItemType
    
    
}

struct JobLineItem : LineItem {
    let hash: Hash
    let creation: Date
    let deadline: Date
    let description: Description
    
    func toString(complete: Bool = false) -> String {
        let createStr = VxdayUtil.datetimeFormatter.string(from: creation)
        let deadlineStr = VxdayUtil.dateFormatter.string(from: deadline)
        let itemStr = complete ? ItemType.complete.rawValue : itemType().rawValue
        return itemStr + " " + hash.hash + " " + createStr + " " + deadlineStr + " " + description.text
    }
    func itemType() -> ItemType {
        return .job
    }
}

struct TaskLineItem : LineItem {
    let hash: Hash
    let description: Description
    
    func toString(complete: Bool = false ) -> String {
        let itemStr = complete ? ItemType.complete.rawValue : itemType().rawValue
        return itemStr + " " + hash.hash + " " + description.text
    }
    func itemType() -> ItemType {
        return .task
    }
}

struct TokenEntry : LineItem {
    let hash: Hash
    let start: Date
    let stop: Date
    let description: Description = Description("")
    func toString(complete: Bool = false) -> String {
        
        let startStr = VxdayUtil.datetimeFormatter.string(from: start)
        let stopStr = VxdayUtil.datetimeFormatter.string(from: stop)
        let humanReadable = VxdayUtil.humanDuration(between: start, and: stop)
        return itemType().rawValue + " " + startStr + " " + stopStr + " " + humanReadable
    }
    
    func itemType() -> ItemType {
        return .tokenEntry
    }
}

class ItemParser {
    
    static func create(from line: String) -> LineItem? {
        let array = VxdayUtil.splitString(line)
        
        guard let itemType = ArgParser.itemType(args: array, index: 0) else {
            print("Error reading the item type from array: \(array). Example: -. 01234567  That hyphen dot denotes valid item type.")
            return nil
        }
    
        guard let hash = ArgParser.hash(args: array, index: 1) else {
            print("Error reading hash from item line: \(array)")
            return nil
        }
        
        switch itemType {
            case .job:
                guard let creationDate = ArgParser.dateTime(args: array, index: 2) else {
                    print("Error: could not extract creation date from: \(array)")
                    return nil
                }
                guard let deadline = ArgParser.date(args: array, index: 3) else {
                    print("Error: could not extract deadline from: \(array)")
                    return nil
                }
                guard let description = ArgParser.description(args: array, start: 4) else {
                    print("Error: could not get description from: \(array)")
                    return nil
                }
                return JobLineItem(hash: hash, creation: creationDate, deadline: deadline, description: description)
            case .task:
                guard let description = ArgParser.description(args: array, start: 2) else {
                    print("Error: could not get description from: \(array)")
                    return nil
                }
                return TaskLineItem(hash: hash, description: description)
            case .tokenEntry:
                print("TODO")
            return nil
        case .complete:
            print("This shouldnt happen.")
            return nil 
            
        }
    }
}


struct Item {
    var type : ItemType
    var created: Date
    var completion: Date
    var description: String
    var hash: String
    init?(_ line: String) {
        print("hash of line is: \(VxdayUtil.hash(line))")
        let array = VxdayUtil.splitString(line)
        guard array.count > 3  else {
            print("Error parsing. : \(line), not enough items in array: \(array)")
            return nil
        }
        guard let type = ItemType(rawValue: array[0]) else {
            print("Error parsing: \(line), unknown Vxday type: \(array[0])")
            return nil
        }
        
        if type == .tokenEntry {
            print("TODO handle token entry")
            return nil
        }
        else {
            let hash = array[1]
            guard VxdayUtil.isValidHash(hash) else {
                print("invalid hash.: \(hash)")
                return nil
            }
            
            guard let createdTime = VxdayUtil.datetimeFormatter.date(from: array[2]) else {
                print("Error extracting created time: \(array[2])")
                return nil
            }
            guard let deadline = VxdayUtil.dateFormatter.date(from: array[3]) else {
                print("Error extracted completion date from : \(array[3])")
                return nil
            }
            guard let description = VxdayUtil.flattenRest(array, start: 4) else {
                print("Error, this job has no description.")
                return nil
            }
            self.hash = hash
            self.type = type
            self.created = createdTime
            self.completion = deadline
            self.description = description
        }
    }
    
    func toString() -> String {
        let typeStr = type.rawValue
        let createdStr = VxdayUtil.datetimeFormatter.string(from: self.created)
        let completionStr = VxdayUtil.dateFormatter.string(from: self.completion)
        return typeStr + " " + createdStr + " " + completionStr + " " + description
    }
}


class VxdayInstruction {
    
    static func makeAddString(_ listName: ListName, description: Description, offset: IntOffset?) -> String {
        
        let now = VxdayUtil.now()
        let created = VxdayUtil.datetimeFormatter.string(from: now)
        
        let hash = VxdayUtil.hash(created + description.text)
        
        if let o = offset {
            
            let deadline = VxdayUtil.dateFormatter.string(from: VxdayUtil.increment(date: now, byDays: o.offset))
            
            return "\(ItemType.job.rawValue) \(hash) \(created) \(deadline) \(description.text)"
        }
        else {
            return " \(ItemType.task.rawValue) \(hash) \(description.text)"
        }
    }
    
}
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
    
    private static let bashDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY2_SRC_DIR")! + "/bash"
    }()
    
    private static let activeDir: String = {
        VxdayExec.getEnvironmentVar("VXDAY2_ACTIVE_DIR")!
    }()
    
    private static let retiredDir: String = {
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
//
//  Vxday_util.swift
//  Vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class VxdayUtil {
    
    private static let datetimeFormat = "yyyy-MM-dd'T'HH:mm:ss"
    private static let dateFormat = "yyyy-MM-dd"
    
    static let datetimeFormatter : DateFormatter  = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = VxdayUtil.datetimeFormat
        return dateFormatter
    }()
    
    static let dateFormatter : DateFormatter = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = VxdayUtil.dateFormat
        return dateFormatter
    }()
    
    class func splitString(_ string: String) -> [String] {
        return string.components(separatedBy: " ")
        
    }
    
    class func isValidHash(_ hash: String) -> Bool {
        return hash.characters.count == 5
    }
    
    
    class func flattenRest(_ array: [String], start: Int) -> String? {
        if array.count < start {
            print("error flattening the rest.")
            return nil
        }
        let endWords = array.suffix(array.count - start)
        return endWords.flatMap({$0 + " " }).joined()
    }
    
    class func now() -> Date {
        return Date()
    }
    
    
    class func hash(_ string: String) -> String {
        let time = now().timeIntervalSince1970.hashValue
        let hashable = "\(time)" + string
        let hashed = "0" + String(format: "%2X", hashable).lowercased()
        return hashed
    }
    
    class func increment(date: Date, byDays days : Int) -> Date {
        var comp = DateComponents()
        comp.day = days
        return  Calendar.current.date(byAdding: comp, to: date )!
    }
    
    class func humanDuration(between start: Date, and end: Date) -> String {
        return "TODO"
        
    }
    
}

//
//  VxdayParser.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation




struct Hash {
    let hash: String
    init(_ hash: String) {
        self.hash = hash
    }
}

struct Description {
    let text: String
    init(_ text: String) {
        self.text = text
    }
}

struct ListName {
    let name: String
    init(_ name: String) {
        self.name = name
    }
}
struct IntOffset {
    let offset: Int
    init(_ offset: Int) {
        self.offset = offset
    }
}




enum Verb : String {
    case add = "add"
    case all = "all"
    case doIt = "do"
    case go = "go"
    case help = "help"
    case less = "less"
    case note = "note"
    case retire = "retire"
    case today = "today"
    case top = "top"
    case track = "track"
    case unretire = "unretire"
    case what = "what"
    case x = "x"
    case jot = "jot"
    
    func isa(_ verbs: [Verb]) -> Bool {
        for v in verbs {
            if v == self {
                return true
            }
        }
        return false
    }
}

enum Instruction {
    
    //list actions
    case retire(ListName)
    case unretire(ListName)
    case add(ListName, IntOffset, Description)
    case doIt(ListName, Description)
    case lessList(ListName)
    case allList(ListName)
    case trackList(ListName)
    case top(ListName)
    
    
    //hash actions
    case x(Hash)
    case go(Hash)
    case note(Hash)
    case lessHash(Hash)
    case trackHash(Hash)
    
    //global actions
    
    case today(IntOffset?)
    case all
    case what
    case track
    case week(IntOffset?)
    case help
}



class ArgParser {
    
    static func createInstruction(args: [String]) -> Instruction? {
        guard args.count > 0 else  {
            print("Error: Empty instruction args. You need a verb. Example: day help")
            return nil
        }
        guard let verb = Verb(rawValue: args[0]) else {
            print("Error: This isn't a valid verb: \(args[0]). Example: day help")
            return nil
        }
        
        switch verb {
        case .add:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Add couldn't find list name in \(args)")
                return nil
            }
            guard let offset = ArgParser.offset(args: args, index: 2) else {
                print("Error: Add couldn't find offset in args: \(args)")
                return nil
            }
            guard let description = ArgParser.description(args: args, start: 3) else {
                print("Error: Add couldn't find a description in args: \(args)")
                return nil
            }
            return .add(listName, offset, description)
            
        case .all:
            if let listName = ArgParser.listName(args: args, index: 1) {
                return .allList(listName)
            }
            else {
                return .all
            }
            
        case .doIt:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Do couldn't find list name in \(args)")
                return nil
            }
            guard let description = ArgParser.description(args: args, start: 2) else {
                print("Error: Do couldn't find a description in args: \(args)")
                return nil
            }
            return .doIt(listName, description)
            
        case .go:
            guard let hash = ArgParser.hash(args: args, index: 1) else {
                print("Error: Couldn't find hash  in \(args)")
                return nil
            }
            return .go(hash)
            
        case .help:
            return .help
            
        case .jot:
            guard let description = ArgParser.description(args: args, start: 1) else {
                print("Error: Do couldn't find a description in args: \(args)")
                return nil
            }
            return .doIt(ListName("_jot"), description)
            
        case .less:
            if let hash = ArgParser.hash(args: args, index: 1)  {
                return .lessHash(hash)
            }
            else {
                guard let listName = ArgParser.listName(args: args, index: 1) else {
                    print("Error: Do couldn't find list name or hash name in \(args)")
                    return nil
                }
                return .lessList(listName)
            }
            
        case .note:
            guard let hash = ArgParser.hash(args: args, index: 1) else {
                print("Error: Couldn't find hash name in \(args)")
                return nil
            }
            return .note(hash)
            
            
        case .retire:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Retire couldn't find list name in \(args)")
                return nil
            }
            return .retire(listName)
            
            
        case .today:
            return .today(ArgParser.offset(args: args, index: 1)) // nil is ok.
            
        case .top:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Do couldn't find list name in \(args)")
                return nil
            }
            return .top(listName)
            
        case .track:
            if let hash = ArgParser.hash(args: args, index: 1)  {
                return .trackHash(hash)
            }
            
            if let listName = ArgParser.listName(args: args, index: 1) {
                return .trackList(listName)
            }
            return .track
            
            
            
        case .unretire:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Unretire couldn't find list name in \(args)")
                return nil
            }
            return .unretire(listName)
            
            
        case .what:
            return .what
            
        case .x:
            guard let hash = ArgParser.hash(args: args, index: 1) else {
                print("Error: Couldn't find hash  in \(args)")
                return nil
            }
            return .x(hash)
        }
    }
    
    static func listName(args: [String], index: Int) -> ListName? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        if str.characters.first == "0" {
            print("Error, this looks like a hash, not a list name: \(str)")
            return nil
        }
        return ListName(str)
    }
    
    static func offset(args: [String], index: Int) -> IntOffset? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        guard let intOffset : Int  = Int(str) else {
            print("Error, this doesn't look like an integer offset: \(str)")
            return nil
        }
        return IntOffset(intOffset)
    }
    
    static func hash(args: [String], index: Int) -> Hash? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        guard  str.characters.first != "0" else {
            print("Error this doesn't look like a hash, it doesn't start with a 0: \(str)")
            return nil
        }
        guard str.characters.count == 8 else {
            print("This hash is the wrong length: \(str). Hashes are 8 chars")
            return nil
        }
        return Hash(str)
        
    }
    
    static func description(args: [String], start: Int) -> Description? {
        guard let str = VxdayUtil.flattenRest(args, start: start) else {
            return nil
        }
        return Description(str)
    }
    static func itemType(args: [String], index: Int) -> ItemType? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return ItemType(rawValue: str)
    }
    
    static func dateTime(args: [String], index: Int) -> Date? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return VxdayUtil.datetimeFormatter.date(from: str)
        
    }
    static func date(args: [String], index: Int) -> Date? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return VxdayUtil.dateFormatter.date(from: str)
    }
    
    private static func str(args: [String], index: Int) -> String? {
        guard args.count > index else {
            return nil
        }
        return args[index]
    }
    
}
//
//  test.swift
//  vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//


import Foundation

func printArgs() {
    for  argument in CommandLine.arguments {
        print(argument)
    }
}




/*
let original = "-. abcaa 2013-05-08T19:03:53 2013-05-08 One Two Three Four."
let item = Item(original)
print("original text is :")
print(original)
print("item is ")
print(item!)
print("converted back its:")
print(original)
printArgs()
*/
//shell("echo", "wehey this works")
//shell("ls", "-a")
//shell("touch", "testfile.txt")
//shell("./vic.sh", "from swift!")
//VxdayExec.retire(ListName("me"))
//VxdayExec.unretire(ListName("me"))




/*
let str = VxdayInstruction.makeAddString(ListName("me"), description: Description("Check add string looks correct."), offset: IntOffset(4))
print("str is \(str)")
let task = VxdayInstruction.makeAddString(ListName("vxday"), description: Description("Get task strings wokring."), offset: nil)
print("task is \(task)")
VxdayExec.append(ListName("wehey"), content: str)
*/

let now = VxdayUtil.now()

//VxdayExec.note(ListName("bam"), hash: Hash("abcdefg"))




print("waiting: \(now)")
VxdayExec.wait(ListName("me"), hash: Hash("abcdefg"))
/*
if let x = readLine(strippingNewline: true) {
    print("Read line \(x)")
}


let finish = VxdayUtil.now()
print("done waiting. \(finish)")

 */

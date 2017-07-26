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

struct CompletionDate {
    let date: Date
    init(_ date: Date) {
        self.date = date
    }
    init?(from string: String) {
        guard let d = VxdayUtil.datetimeFormatter.date(from: string) else {
            return nil
        }
        self.date = d
    }
    
    func pretty() -> String {
        return self.date.ago()
    }
    func toString() -> String {
        return VxdayUtil.datetimeFormatter.string(from: self.date)
    }
}
struct CreationDate {
    let date: Date
    init(_ date: Date) {
        self.date = date
    }
    init?(from string: String) {
        guard let d = VxdayUtil.datetimeFormatter.date(from: string) else {
            return nil
        }
        self.date = d
    }
    
    func pretty() -> String {
        return self.date.ago()
    }
    
    func toString() -> String {
        return VxdayUtil.datetimeFormatter.string(from: self.date)
    }
}

struct DeadlineDate {
    let date: Date
    init(_ date: Date) {
        self.date = date
    }
    init?(from string: String) {
        guard let d = VxdayUtil.dateFormatter.date(from: string) else {
            return nil
        }
        self.date = d
    }
    func pretty() -> String {
        
        return self.date.daysAgo() 
    }
    
    func toString() -> String {
        return VxdayUtil.dateFormatter.string(from: self.date)
    }
}

struct ListName : Hashable {
    let name: String
    init(_ name: String) {
        self.name = name
    }
    
    var hashValue: Int {
        return name.hashValue
    }
    
    static func == (lhs: ListName, rhs: ListName) -> Bool {
        return lhs.name == rhs.name
    }
}
struct IntOffset {
    let offset: Int
    init(_ offset: Int) {
        self.offset = offset
    }
}


extension Date {
    static func daysOffsetString(_ days: Int) -> String {
        switch days {
            case let x where x == -1:
                return "Yesterday"
            case let x where x == 0:
                return "Today"
            case let x where x == 1:
                return "Tomorrow"
            case let x where x < 0:
                return "\(abs(days)) days ago"
            default:
                return "In \(days) days"
        }
    }
    
    func daysAgoInt() -> Int {
        return Int(self.timeIntervalSince(VxdayUtil.nowDay())) / 86400
    }
    
    func daysAgo() -> String {
        return Date.daysOffsetString(self.daysAgoInt())
    }
    
    func ago() -> String {
        let SECONDS_IN_A_DAY = 86400
        let SECONDS_IN_AN_HOUR = 60 * 60
        let SECONDS_IN_A_MINUTE = 60
        let now = VxdayUtil.now()
        

        let interval =  Int(self.timeIntervalSince(now))
        print("interval is \(interval)")
        let absInterval = abs(interval)
        if absInterval > SECONDS_IN_A_DAY {
            return Date.daysOffsetString(interval / SECONDS_IN_A_DAY)
        }
        else if absInterval > SECONDS_IN_AN_HOUR {
            let hours = interval / SECONDS_IN_AN_HOUR
            if hours < 0 {
                return "In \(hours) hours."
            }
            
            return "\(hours) hours ago"
        }
        else if absInterval > SECONDS_IN_A_MINUTE {
            let mins = interval / SECONDS_IN_A_MINUTE
            if mins < 0 {
                return "In \(mins) minutes."
            }
            return "\(mins) minutes ago."
        }
        return "Just now."
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
    case retired = "retired"
    case today = "today"
    case top = "top"
    case track = "track"
    case unretire = "unretire"
    case what = "what"
    case x = "x"
    case jot = "jot"
    
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
    case retired
    case what
    case track
    case week(IntOffset?)
    case help
    
    
    static func create(_ args:[String]) -> Instruction? {
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
                
            case .retired:
                return .retired
                
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
}



class ArgParser {
    
  
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
        guard  str.characters.first == "0" else {
            print("Error this doesn't look like a hash, it doesn't start with a 0: \(str)")
            return nil
        }
        
        //06b77a160 valid hash.
        guard str.characters.count == 9 else {
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
    
    static func creation(args: [String], index: Int) -> CreationDate? {
        guard let d = ArgParser.dateTime(args: args, index: index) else  {
            return nil
        }
        return CreationDate(d)
    }
    static func completion(args: [String], index: Int) -> CompletionDate? {
        guard let d = ArgParser.dateTime(args: args, index: index) else  {
            return nil
        }
        return CompletionDate(d)
    }
    static func deadline(args: [String], index: Int) -> DeadlineDate? {
        guard let d = ArgParser.date(args: args, index: index) else  {
            return nil
        }
        return DeadlineDate(d)
    }
    
    private static func dateTime(args: [String], index: Int) -> Date? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return VxdayUtil.datetimeFormatter.date(from: str)
    }
    
    private  static func date(args: [String], index: Int) -> Date? {
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
//  VxdayView.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation



protocol VxItem {
    var list: ListName {get}
    var hash: Hash {get}
    var creation: CreationDate {get}
}

struct VxJob : VxItem {
    let list: ListName
    let hash: Hash
    let creation: CreationDate
    let deadline: DeadlineDate
    let description : Description
    let completion: CompletionDate?
    
    func isComplete() -> Bool {
        return completion != nil
    }
}
struct VxTask : VxItem {
    let list: ListName
    let hash: Hash
    let creation: CreationDate
    let description: Description
    let completion : CompletionDate?
    
    func isComplete() -> Bool {
        return completion != nil
    }
}

struct VxToken : VxItem {
    let list: ListName
    let hash: Hash
    let creation: CreationDate
    let completion: CompletionDate
}

enum Item {
    
    
    case job(VxJob)
    case task(VxTask)
    case token(VxToken)
    
    func getJob() -> VxJob? {
        if case let Item.job(job) = self  {
            return job
        }
        return nil
    }
    func getTask() -> VxTask? {
        if case let Item.task(task) = self {
            return task
        }
        return nil
    }
    func getToken() -> VxToken? {
        if case let Item.token(token) = self  {
            return token
        }
        return nil
    }
    
    func vxItem() -> VxItem {
        switch self {
        case let .job(job):
            return job
        case let .task(task):
            return task
        case let .token(token):
            return token
        }
    }
    
    func toString() -> String {
        switch self {
            case let .job(job):
                var completion = job.completion?.toString()  ?? ""
                if completion != "" {
                    completion = completion + " "
                }
                return "\(ItemType.job.rawValue) \(job.hash.hash) \(job.creation.toString()) \(job.deadline.toString()) \(completion)\(job.description.text)"
            case let .task(task):
                var completion = task.completion?.toString()  ?? ""
                if completion != "" {
                    completion = completion + " "
                }
                return "\(ItemType.task.rawValue) \(task.hash.hash) \(task.creation.toString()) \(completion)\(task.description.text)"
            case let .token(token):
                return "\(ItemType.tokenEntry.rawValue) \(token.hash.hash) \(token.creation.toString()) \(token.completion.toString())"
        }
    }
    
    static func create(_ line: String, list: ListName) -> Item? {

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
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let deadlineDate = ArgParser.deadline(args: array, index: 3) else {
                print("Error: could not extract deadline from: \(array)")
                return nil
            }
            guard let description = ArgParser.description(args: array, start: 4) else {
                print("Error: could not get description from: \(array)")
                return nil
            }
            return Item.job(VxJob(list: list, hash: hash, creation: creationDate, deadline: deadlineDate, description: description, completion: nil))
            
        case .task:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            
            guard let description = ArgParser.description(args: array, start: 3) else {
                print("Error: could not get description from: \(array)")
                return nil
            }
            return Item.task(VxTask(list: list, hash: hash, creation: creationDate, description: description, completion: nil ))
        case .tokenEntry:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let completionDate = ArgParser.completion(args: array, index: 2) else {
                print("Error: could not extract completion date from: \(array)")
                return nil
            }
            return Item.token(VxToken(list: list, hash: hash, creation: creationDate, completion: completionDate))
        default:
            print("TODO unhandled vxday line: \(line).")
            return nil
        }
    }
}



class VxdayColor {
    static let dangerColor : String = ANSIColors.red.rawValue
    static let warningColor : String = ANSIColors.yellow.rawValue
    static let baseColor : String = ANSIColors.white.rawValue
    static let happyColor: String = ANSIColors.green.rawValue
    static let titleColor :String = ANSIColors.cyan.rawValue
    
    static func danger(_ string: String) -> String {
        return dangerColor + string + baseColor
    }
    static func warning(_ string: String) -> String {
        return warningColor + string + baseColor
    }
    static func title(_ string: String) -> String {
        return titleColor + string + baseColor
    }
    static func happy(_ string: String) -> String {
        return happyColor + string + baseColor
    }
}

enum ANSIColors: String {
    case black = "\u{001B}[0;30m"
    case red = "\u{001B}[0;31m"
    case green = "\u{001B}[0;32m"
    case yellow = "\u{001B}[0;33m"
    case blue = "\u{001B}[0;34m"
    case magenta = "\u{001B}[0;35m"
    case cyan = "\u{001B}[0;36m"
    case white = "\u{001B}[0;37m"
    
    func name() -> String {
        switch self {
        case .black: return "Black"
        case .red: return "Red"
        case .green: return "Green"
        case .yellow: return "Yellow"
        case .blue: return "Blue"
        case .magenta: return "Magenta"
        case .cyan: return "Cyan"
        case .white: return "White"
        }
    }
    
    static func all() -> [ANSIColors] {
        return [.black, .red, .green, .yellow, .blue, .magenta, .cyan, .white]
    }
}

class VxdayView {
    
    let items: [Item]
    
    init(_ items: [Item]) {
        self.items = items
    }
    
    private func getDeadlines() -> [VxJob] {
        let jobs : [VxJob] = items.flatMap {
            if case let Item.job(job) = $0 {
               return job
            }
            return nil
            }.sorted { $0.deadline.date < $1.deadline.date}
        return jobs
    }
    
    private func getTasks() -> [VxTask] {
        let tasks: [VxTask] = items.flatMap {
            if case let Item.task(task) = $0 {
                return task
            }
            return nil
            }.sorted {$0.creation.date < $1.creation.date}
        return tasks
    }
    
    private func allLists() -> [ListName] {
        var set : Set<ListName> = Set()
        items.map { $0.vxItem().list}.forEach { set.insert($0) }
        return Array(set)
    }

    func showTasks() -> [String] {
        return self.getTasks().map { "Created :" + $0.creation.date.daysAgo() + " : " + $0.description.text }
    }
    
    private func pad(_ string: String, toLength length: Int) -> String {
        if string.characters.count > length {
            return string
        }
        var spaces = ""
        let needed = length - string.characters.count
        for _ in 1...needed {
            spaces = spaces +  " "
        }
        return string +  spaces
    }
    
    func showJobs() -> [String] {
        
        return self.getDeadlines().map {
            
            let daysAgo = $0.deadline.date.daysAgoInt()
            let  timeStr =  pad( $0.deadline.pretty(), toLength: 18)
            var datedStr = timeStr +  VxdayUtil.dateFormatter.string(from: $0.deadline.date) + "   "
            if daysAgo < 0 {
                datedStr = VxdayColor.danger(datedStr)
            }
            if daysAgo == 0 {
                datedStr = VxdayColor.warning(datedStr)
            }
            else {
                datedStr = VxdayColor.happy(datedStr)
            }
            
            return datedStr +  $0.description.text
            
        }
    }
    
    func renderAll() -> [String] {
        var output : [String] = []
        let lists = self.allLists()
        if lists.count == 0 {
            return []
        }
        if lists.count == 1 {
            output.append("Summary for \(lists[0].name):")
        }
      //  let white = ANSIColors.white.rawValue
        output.append("")
        output.append( ANSIColors.green.rawValue + "---------- Tasks ----------")
        output += self.showTasks()
        output.append("")
        output.append( ANSIColors.yellow.rawValue + "----------- Jobs ----------")
        //output.append(ANSIColors.green.rawValue)
        output += self.showJobs()
        output.append(ANSIColors.white.rawValue)
        return output
    }
}

/*
protocol LineItem {
    var list: ListName {get}
    var hash: Hash {get}
    var creation: Date {get}
    var description: Description {get}
    func toString(complete: Bool) -> String
}

struct JobLineItem : LineItem {
    let list: ListName 
    let hash: Hash
    let creation: Date
    let deadline: Date
    let description: Description
    
    func toString(complete: Bool = false) -> String {
        let createStr = VxdayUtil.datetimeFormatter.string(from: creation)
        let deadlineStr = VxdayUtil.dateFormatter.string(from: deadline)
        let itemStr = complete ? ItemType.complete.rawValue : ItemType.job.rawValue
        return itemStr + " " + hash.hash + " " + createStr + " " + deadlineStr + " " + description.text
    }
}

struct TaskLineItem : LineItem {
    let list: ListName
    let hash: Hash
    let description: Description
    let creation: Date
    func toString(complete: Bool = false ) -> String {
        let itemStr = complete ? ItemType.complete.rawValue : ItemType.task.rawValue
        return itemStr + " " + hash.hash + " " + description.text
    }
    func itemType() -> ItemType {
        return .task
    }
}

struct TokenEntry : LineItem {
    let list: ListName
    let hash: Hash
    let creation: Date
    let stop: Date
    let description: Description = Description("")
    func toString(complete: Bool = false) -> String {
        
        let startStr = VxdayUtil.datetimeFormatter.string(from: creation)
        let stopStr = VxdayUtil.datetimeFormatter.string(from: stop)
        let humanReadable = VxdayUtil.humanDuration(between: creation, and: stop)
        return ItemType.tokenEntry.rawValue + " " + startStr + " " + stopStr + " " + humanReadable
    }
    
}
*/



//
//  VxdayInstruction.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

class VxdayInstruction {
    
    static func executeInstruction(_ instruction : Instruction) {
        switch instruction {
            case let .add(list, offset, description):
                let string = VxdayInstruction.makeAddString(description: description, offset: offset)
                VxdayExec.append(list, content: string)
            
            case let .doIt(list, description):
                let string = VxdayInstruction.makeAddString(description: description , offset: nil)
                VxdayExec.append(list, content: string)
            case let .retire(list):
                VxdayExec.retire(list)
            case let .unretire(list):
                VxdayExec.unretire(list)
            case let .lessList(list):
                VxdayExec.lessList(list)
            case let .allList(list):
                VxdayExec.allList(list)
           // case let .x(hash):
           //     VxdayExec.x(hash)
            
            
        default:
             print("TODO handle instruction: \(instruction)")
        }
    }

    static func makeAddString(description: Description, offset: IntOffset?) -> String {
        
        let now = VxdayUtil.now()
        let created = VxdayUtil.datetimeFormatter.string(from: now)
        
        let hash = VxdayUtil.hash(created + description.text)
        
        if let o = offset {
            let deadline = VxdayUtil.dateFormatter.string(from: VxdayUtil.increment(date: now, byDays: o.offset))
            
            return "\(ItemType.job.rawValue) \(hash) \(created) \(deadline) \(description.text)"
        }
        else {
            return " \(ItemType.task.rawValue) \(hash) \(created) \(description.text)"
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
    static func allList(_ list: ListName) {
        let filename = VxdayFile.getSummaryFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.readSummary(contents, list: list)
        let view  = VxdayView(items)
        let jobsStrings = view.renderAll()
        jobsStrings.forEach { print($0)}
    }
    
    //TODO try to write these using mv
    
    static func red() {
        VxdayExec.shell("tput", "sgr", "0")
    }
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
    
    class func flattenRest(_ array: [String], start: Int) -> String? {
        if array.count < start {
            print("error flattening the rest.")
            return nil
        }
        let endWords = array.suffix(array.count - start)
        return endWords.flatMap({$0 + " " }).joined()
    }
    
    
    class func nowDay() -> Date {
        let c = Calendar.current
        return c.startOfDay(for: now())
    }
    
    class func now() -> Date {
        return Date()
    }
    
    
    class func hash(_ string: String) -> String {
        let time = now().timeIntervalSince1970.hashValue
        let hashable = "\(time)" + string
        var hashed = "0" + String(format: "%2X", hashable).lowercased()
        if hashed.characters.count < 9 {
            let numNeeded = 9 - hashed.characters.count
            for _ in 1...numNeeded {
                hashed += "0"
            }
        }
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
    
    class func beforeUnderscore(_ string: String) -> String {
        return string.components(separatedBy: "_").first ?? ""
    }
    
}

//
//  VxdayRead.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

//TODO make this config
enum ItemType : String {
    case complete = "x."
    case completeJob = "X."
    case tokenEntry = "->."
    case job = "=."
    case task = "-."
}

class VxdayReader {
    
    static func allLists() -> [ListName] {
        let fm = FileManager.default
        let enumerator = fm.enumerator(atPath: VxdayFile.activeDir)!
        var lists: Set<String> = Set()
        while  let file  = enumerator.nextObject() as? String  {
            lists.insert(VxdayUtil.beforeUnderscore(file))
        }
        return lists.map {return ListName($0)}
    }
    
    static func readFile(_ path: String) -> [String] {
        guard let contents =  try? String(contentsOfFile: path) else {
            print("Error reading file: \(path)")
            return []
        }
        return contents.components(separatedBy: "\n").filter{ $0 != ""}
    }
    
    
    static func readSummary(_ lines: [String], list: ListName) -> [Item] {
        return lines.flatMap{ Item.create($0, list: list)}
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



//VxdayExec.note(ListName("bam"), hash: Hash("abcdefg"))



//print("waiting: \(now)")
//VxdayExec.wait(ListName("me"), hash: Hash("abcdefg"))
/*
if let x = readLine(strippingNewline: true) {
    print("Read line \(x)")
}


let finish = VxdayUtil.now()
print("done waiting. \(finish)")

 */

/*
 let location = "/Users/vic/Desktop/test.txt"
 let x =  try? String(contentsOfFile: location)
 print("x is \(x!)")

 */

/*
let allLists = VxdayReader.allLists()
print("all Lists: \(allLists)")
*/

let list = ListName("vic")
VxdayExec.allList(list)

/*
let summaryPath = VxdayFile.getSummaryFilename(list)
let contents = VxdayReader.readFile(summaryPath)
print("CONTENTS ARE: \(contents)'")
let items = VxdayReader.readSummary(contents, list: list)
print("items are: \(items)")
 */




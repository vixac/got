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
    
    func isValid() -> Bool {
        return hash.characters.first == "0" && hash.characters.count == 9
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
    case complete = "complete"
    case remove = "remove"
    
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
    case complete(ListName?)
    
    //hash actions
    case x(Hash)
    case go(Hash)
    case note(Hash)
    case lessHash(Hash)
    case trackHash(Hash)
    case remove(Hash)
    
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
            
            case .remove:
                guard let hash = ArgParser.hash(args: args, index: 1) else {
                    print("Error: Couldn't find hash name in \(args)")
                    return nil
                }
                return .remove(hash)
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
                
            
            case .complete:
                if let listName = ArgParser.listName(args: args, index: 1) {
                    return .complete(listName)
                }
                else {
                    return .complete(nil)
                }
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
    func toVxday() -> String
    func itemType() -> ItemType
    func complete() -> VxItem
    func isComplete() -> Bool
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
    
    
    func complete()  -> VxItem {
        return VxJob(list: list, hash: hash, creation: creation, deadline: deadline, description: description, completion: CompletionDate(VxdayUtil.now()))
    }
    
    func toVxday() -> String {
        let symbolStr = itemType().rawValue + " "
        let hashStr = hash.hash + " "
        let creationStr = creation.toString() + " "
        let deadlineStr = deadline.toString() + " "
        let descriptionStr = description.text + " "
        let completionStr = completion != nil  ? (VxdayUtil.datetimeFormatter.string(from: completion!.date) + " ") : ""
        let str = symbolStr + hashStr + creationStr + deadlineStr + completionStr + descriptionStr
        return str
    }
    
    func itemType() -> ItemType {
        if isComplete() {
            return .completeJob
        }
        return .job
    }
    
}
struct VxTask : VxItem {
    let list: ListName
    let hash: Hash
    let creation: CreationDate
    let description: Description
    let completion : CompletionDate?
    
    func complete() -> VxItem {
        return VxTask(list: list, hash: hash, creation: creation, description: description, completion: CompletionDate(VxdayUtil.now()))
    }
    
    func isComplete() -> Bool {
        return completion != nil
    }
    func toVxday() -> String  {
        let symbolStr = itemType().rawValue + " "
        let hashStr = hash.hash + " "
        let creationStr = creation.toString() + " "
        let descriptionStr = description.text + " "
        let completionStr = completion != nil ? (VxdayUtil.datetimeFormatter.string(from: completion!.date) + " ") : ""
        return symbolStr + hashStr + creationStr + completionStr +  descriptionStr
    }
    func itemType() -> ItemType {
        if isComplete() {
            return .completeTask
        }
        return .task
    }
}

struct VxToken : VxItem {
    let list: ListName
    let hash: Hash
    let creation: CreationDate
    let completion: CompletionDate
    
    func complete()  -> VxItem{
        return self
    }
    
    func toVxday() -> String  {
        return "TODO"
    }
    
    func isComplete() -> Bool {
        return true
    }
    func itemType() -> ItemType {
        return .token
    }
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
                return "\(ItemType.token.rawValue) \(token.hash.hash) \(token.creation.toString()) \(token.completion.toString())"
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
        case .token:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let completionDate = ArgParser.completion(args: array, index: 2) else {
                print("Error: could not extract completion date from: \(array)")
                return nil
            }
            return Item.token(VxToken(list: list, hash: hash, creation: creationDate, completion: completionDate))
            
        case .completeJob:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let deadlineDate = ArgParser.deadline(args: array, index: 3) else {
                print("Error: could not extract deadline from: \(array)")
                return nil
            }
            guard let completionDate = ArgParser.completion(args: array, index: 4) else {
                print("Error: could not extract completion date from : \(array)")
                return nil
            }
            guard let description = ArgParser.description(args: array, start: 5) else {
                print("Error: could not get description from: \(array)")
                return nil
            }
            return Item.job(VxJob(list: list, hash: hash, creation: creationDate, deadline: deadlineDate, description: description, completion: completionDate))
        case .completeTask:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let completionDate = ArgParser.completion(args: array, index: 3) else {
                print("Error: could not extract completion date from : \(array)")
                return nil
            }
            guard let description = ArgParser.description(args: array, start: 4) else {
                print("Error: could not get description from: \(array)")
                return nil
            }
            return Item.task(VxTask(list: list, hash: hash, creation: creationDate, description: description, completion: completionDate))
        }
    }
}



class VxdayColor {
    static let dangerColor : String = ANSIColors.red.rawValue
    static let warningColor : String = ANSIColors.yellow.rawValue
    static let baseColor : String = ANSIColors.reset.rawValue
    static let happyColor: String = ANSIColors.green.rawValue
    static let titleColor :String = ANSIColors.white.rawValue
    static let info2Color: String = ANSIColors.test.rawValue
    static let infoColor: String = ANSIColors.cyan.rawValue
    static let whiteColor : String = ANSIColors.white.rawValue
    
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
    static func info(_ string: String) -> String {
        return infoColor + string + baseColor
    }
    static func info2(_ string: String) -> String {
        return info2Color + string + baseColor
    }
    static func boldInfo(_ string: String) -> String {
        return whiteColor + string + baseColor
    }
    
    static func putBack() {
        print(ANSIColors.reset.rawValue)
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
    case white = "\u{001B}[1;37m"
    case test = "\u{001B}[1;34m"
    case reset = "\u{001B}[0;0m"
}


class Spaces {
    static let List = 9
    static let Timeliness = 14
    static let WhatOverdue = 13
    static let WhatPresent = 13
    static let WhatFuture = 13
    static let WhatTasks = 12
    static let DaysString = 15
    static let Hash = 11
    
}

class ListSummary {
    let list: ListName
    var past: [VxJob] = []
    var present: [VxJob] = []
    var future: [VxJob] = []
    var taskCount: Int = 0
    init(_ list: ListName) {
        self.list = list
    }
    func addJob(_ job: VxJob) {
        let bucket = job.deadline.date.bucket()
        switch bucket {
        case .past:
            past.append(job)
        case .present:
            present.append(job)
        case .future:
            future.append(job)
        }
    }
     func addTask(_ task: VxTask) {
        taskCount += 1
    }
    func total() -> Int {
        return past.count + present.count + future.count + taskCount
    }
}

class VxdayView {
    
    let items: [Item]
    
    init(_ items: [Item]) {
        self.items = items
    }
    
    func toBuckets() -> [ListSummary] {
        var dict : [ListName: ListSummary] = [:]
        for item in items {
            if let job = item.getJob() {
                let list = job.list
                if dict[list] == nil {
                     let summary = ListSummary(list)
                    summary.addJob(job)
                    dict[list] = summary
                }
                else  {
                    dict[list]?.addJob(job)
                }
                
            }
            else if let task = item.getTask() {
                let list = task.list
                if dict[task.list] == nil {
                    let summary =  ListSummary(list)
                    summary.addTask(task)
                    dict[list] = summary
                }
                else {
                    dict[list]?.addTask(task)
                }
                
            }
            //tokens aren't part of the summary yet.
        }
        return dict.map { return $0.value }.sorted( by: { $0.total() < $1.total()})
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
        return self.getTasks().map {
            let dateStr = daysView($0.creation.date)
            let hash = hashView($0.hash)
            let listName = listNameView($0.list)
            return dateStr + hash + listName +  $0.description.text
        }
    }
    
    func space() -> String {
        return "   "
    }
    
    func hashView(_ hash: Hash) -> String {
        return VxdayColor.info2(pad(hash.hash, toLength: Spaces.Hash))
    }
    
    func listNameView(_ list: ListName) -> String {
        return VxdayColor.info(pad(list.name, toLength: Spaces.List))
    }
    

    func timeBucketToColor(_ date: Date, string: String)  -> String {
        let bucket = date.bucket()
        switch bucket {
            case .past:
                return VxdayColor.danger(string)
            case .present:
                return VxdayColor.warning(string)
            case .future:
                return  VxdayColor.happy(string)
            }
    }
    
    func daysView(_ date: Date) -> String {
        let agoStr = pad(date.daysAgo(), toLength: 14)
        let datedStr = agoStr + VxdayUtil.dateFormatter.string(from: date) + "   "
        return timeBucketToColor(date, string: datedStr)
    }

    func showJobs() -> [String] {
        return self.getDeadlines().map {
            let timeStr =  pad( $0.deadline.pretty(), toLength: Spaces.DaysString)
            var datedStr = timeStr +  VxdayUtil.dateFormatter.string(from: $0.deadline.date) + space()
            datedStr = timeBucketToColor($0.deadline.date, string: datedStr)
            
            let hash = hashView($0.hash)
            let listName = listNameView($0.list)
            return daysView($0.deadline.date) + hash + listName +  $0.description.text
            
        }
    }
    
    
    func noStringForZero(_ prefix: String, number: Int, toLength : Int) -> String {
        if number == 0 {
            return pad("", toLength: toLength)
        }
        let str = "\(prefix) \(number)"
        return pad(str, toLength: toLength)
    }
    func globalOneLiner(buckets: [ListSummary]) -> String {
        let listCount = buckets.count
        let overdueCount = buckets.map { $0.past.count}.reduce(0, { $0 + $1})
        let todayCount = buckets.map { $0.present.count}.reduce(0, { $0 + $1})
        let futureCount = buckets.map { $0.future.count}.reduce(0 , {$0 + $1})
        let taskCount = buckets.map {$0.taskCount}.reduce(0, {$0 + $1})
        let totalLists = VxdayColor.boldInfo(pad("Lists: \(listCount)", toLength: Spaces.List))
        let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: overdueCount, toLength:  Spaces.WhatOverdue))
        let present = VxdayColor.warning(noStringForZero("Today:", number: todayCount, toLength: Spaces.WhatPresent))
        let future = VxdayColor.happy(noStringForZero("Upcoming:", number: futureCount, toLength: Spaces.WhatFuture))
        let tasks  = VxdayColor.warning(noStringForZero("Tasks:", number: taskCount, toLength: Spaces.WhatTasks))
        let total = VxdayColor.boldInfo("Total: \(overdueCount + todayCount + futureCount + taskCount)")
        return "\(totalLists) \(overdue) \(present) \(future) \(tasks) \(total)"
    }
    
    
    func oneLiners(_ buckets: [ListSummary]) -> [String] {
        return buckets.map  { summary in
            let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: summary.past.count, toLength: Spaces.WhatOverdue))
            let present = VxdayColor.warning(noStringForZero("Today:", number: summary.present.count, toLength: Spaces.WhatPresent))
            let upcoming = VxdayColor.happy(noStringForZero("Upcoming:", number: summary.future.count, toLength: Spaces.WhatFuture))
            let tasks = VxdayColor.warning(noStringForZero("Tasks:", number: summary.taskCount, toLength: Spaces.WhatTasks))
            let total = "Total: \(summary.past.count + summary.present.count + summary.future.count + summary.taskCount)"
            let listName = listNameView(summary.list)
            return "\(listName) \(overdue) \(present) \(upcoming) \(tasks) \(total)"
        }
    }
    
    func renderComplete() -> [String] {
        var strings: [String] = []
        items.forEach { item in
            var completed : CompletionDate? = nil
            var description : Description = Description("")
            var offsetFromDeadline : Int? = nil
            if let job = item.getJob() {
                completed = job.completion
                
                description = job.description
                if let c = completed {
                    offsetFromDeadline = VxdayUtil.timeliness(deadline: job.deadline, completion: c)
                }
            }
            else if let task = item.getTask() {
                completed = task.completion
                description = task.description
            }
            
            if let c = completed {
                //let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: overdueCount, toLength:  Spaces.WhatOverdue))
                let completed = VxdayColor.boldInfo( pad(c.pretty(), toLength: Spaces.DaysString))
                let d = description.text
                
                var timelinessStr = ""
                if let o = offsetFromDeadline {
                    timelinessStr = VxdayUtil.timelinessToString(offsetFromDeadline!)
                    timelinessStr = pad(timelinessStr, toLength:  Spaces.Timeliness)
                    if o < 0 {
                        timelinessStr = VxdayColor.danger(timelinessStr)
                    }
                    else if o == 0 {
                        timelinessStr = VxdayColor.warning(timelinessStr)
                    }
                    else {
                        timelinessStr = VxdayColor.happy(timelinessStr)
                    }
                }
                else {
                    timelinessStr =  pad("", toLength: Spaces.Timeliness)
                }
                
                let hash = hashView(item.vxItem().hash)
                let list = listNameView(item.vxItem().list)
                strings.append("\(completed) \(timelinessStr) \(hash) \(list) \(d)")
            }
        }
        return strings
    }
    
    func renderAll() -> [String] {
        var output : [String] = []
        let lists = self.allLists()
        if lists.count == 0 {
            return []
        }
        
        output.append("")
        output.append(VxdayColor.title("--------------- Tasks --------------"))
        output += self.showTasks()
        output.append("")
        output.append(VxdayColor.title("--------------- Jobs  --------------"))
        output += self.showJobs()
        
        return output
    }
    
    private func pad(_ string: String, toLength length: Int) -> String {
        if string.characters.count > length {
            return string
        }
        var spaces = ""
        let needed = length - string.characters.count
        if needed < 1 {
            return string
        }
        for _ in 1...needed {
            spaces = spaces +  " "
        }
        return string +  spaces
    }
}




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
                let now = VxdayUtil.now()
                let created = CreationDate(now)
                let deadline = DeadlineDate(VxdayUtil.increment(date: now, byDays: offset.offset))
                let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
                let item = VxJob(list: list, hash: hash, creation: created, deadline: deadline, description: description , completion: nil)
                
                VxdayExec.storeItem(item)
            
            
            case let .doIt(list, description):
                let now = VxdayUtil.now()
                let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
                let created = CreationDate(now)
                let item = VxTask(list: list, hash: hash, creation: created, description: description, completion: nil)
                VxdayExec.storeItem(item)
            case let .retire(list):
                VxdayExec.retire(list)
            case let .unretire(list):
                VxdayExec.unretire(list)
            case let .lessList(list):
                VxdayExec.lessList(list)
            case let .allList(list):
                VxdayExec.allList(list)
            case .all:
                VxdayExec.all()
        case let .complete(list):
            VxdayExec.showComplete(list)
            case .what:
                VxdayExec.what()
            case let .x(hash):
                VxdayExec.x(hash)
        case let .remove(hash):
                VxdayExec.remove(hash)
            
           // case let .x(hash):
           //     VxdayExec.x(hash)
            
            
        default:
             print("TODO handle instruction: \(instruction)")
        }
    }

    
    //TODO RM
    /*
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
    */
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
    case note = "note.sh"
    case wait = "wait.sh"
    case lookForHash = "look_for_hash.sh"
    case removeLine = "remove_line.sh"
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
    
    static let outputFile : String = {
        VxdayExec.getEnvironmentVar("VXDAY2_OUTPUT_FILE")!
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
    
    static func getCompleteFilename(_ list: ListName) -> String  {
        return VxdayFile.activeDir + "/" + list.name + "_complete.vxday"
    }
    static func getTokenFilename(_ list: ListName) -> String {
        return VxdayFile.activeDir + "/" + list.name + "_tokens.vxday"
    }

}

class VxdayExec {
    
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
    
    static func getEnvironmentVar(_ name: String) -> String? {
        guard let rawValue = getenv(name) else { return nil }
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
        let view = VxdayView(allItemsEver)
        let buckets = view.toBuckets()
        
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
    
    static func showComplete(_ list: ListName?) {
        if let l = list {
            let completeFile = VxdayFile.getCompleteFilename(l)
            let contents = VxdayReader.readFile(completeFile)
            let items = VxdayReader.linesToItems(contents, list: l)
            
            let view = VxdayView(items)
            let lines = view.renderComplete()
            lines.forEach { print($0)}
        }
        else {
            var allCompleteItems: [Item] = []
            
            let all = VxdayReader.allLists()
            all.forEach { list in
                allCompleteItems += VxdayReader.completeItemsInList(list)
            }
            let view = VxdayView(allCompleteItems)
            let lines = view.renderComplete()
            lines.forEach { print($0)}
        }
    }
    
    static func remove(_ hash: Hash) {
        guard let list = hashToListName(hash) else {
            return
        }
        removeActiveHash(hash, list: list)
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
    
    
    static func storeItem(_ item: VxItem) {
        let script = VxdayFile.getScriptPath(.append)
        let list  = item.list
        var filename : String = ""
        if item.itemType() == ItemType.token {
            print("TODO write this token somewhere.")
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
    
    enum TimeBucket {
        case past
        case present
        case future
    }
    
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
    
    class func hash(_ string: String) -> Hash {
        let time = now().timeIntervalSince1970.hashValue
        let hashable = "\(time)" + string
        var hashed = "0" + String(format: "%2X", hashable).lowercased()
        if hashed.characters.count < 9 {
            let numNeeded = 9 - hashed.characters.count
            for _ in 1...numNeeded {
                hashed += "0"
            }
        }
        return Hash(hashed)
    }
    
    class func timelinessToString(_ daysOffset: Int) -> String {
        if daysOffset == -1 {
            return "1 day late"
        }
        if daysOffset < 0 {
            return "\(abs(daysOffset)) days late"
        }
        if daysOffset == 0 {
            return "On time"
        }
        if daysOffset == 1 {
            return "1 day early"
        }
        return "\(daysOffset) days early"
    }
    class func timeliness(deadline: DeadlineDate, completion: CompletionDate) -> Int {
        return  deadline.date.daysSince(completion.date)
        
    }
    class func increment(date: Date, byDays days : Int) -> Date {
        var comp = DateComponents()
        comp.day = days
        return  Calendar.current.date(byAdding: comp, to: date )!
    }

    
    class func afterHyphen(_ string: String) -> String {
        return string.components(separatedBy: "/").last ?? string
    }
    
    class func beforeUnderscore(_ string: String) -> String? {
        return string.components(separatedBy: "_").first 
    }
}

extension Date {
    static let SECONDS_IN_A_DAY = 86400
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
    func daysSince(_ date: Date) -> Int {
        return Int(self.timeIntervalSince(date)) / Date.SECONDS_IN_A_DAY
    }
    func daysAgoInt() -> Int {
        return self.daysSince(VxdayUtil.now())
    }
    
    func daysAgo() -> String {
        return Date.daysOffsetString(self.daysAgoInt())
    }
    
    func bucket() -> VxdayUtil.TimeBucket {
        let daysAgo = self.daysAgoInt()
        if daysAgo < 0 {
            return .past
        }
        if daysAgo == 0 {
            return .present
        }
        return .future
    }
    
    func ago() -> String {
        let SECONDS_IN_A_DAY = 86400
        let SECONDS_IN_AN_HOUR = 60 * 60
        let SECONDS_IN_A_MINUTE = 60
        let now = VxdayUtil.now()
        
        
        let interval =  Int(self.timeIntervalSince(now))
        let absInterval = abs(interval)
        if absInterval > SECONDS_IN_A_DAY {
            return Date.daysOffsetString(interval / SECONDS_IN_A_DAY)
        }
        else if absInterval > SECONDS_IN_AN_HOUR {
            let hours = interval / SECONDS_IN_AN_HOUR
            if hours > 0 {
                return "In \(hours) hours."
            }
            
            return "\(abs(hours)) hours ago"
        }
        else if absInterval > SECONDS_IN_A_MINUTE {
            let mins = interval / SECONDS_IN_A_MINUTE
            if mins > 0 {
                return "In \(mins) mins."
            }
            return "\(abs(mins)) mins ago."
        }
        return "Just now."
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
    case completeTask = "x."
    case completeJob = "X."
    case token = "->."
    case job = "=."
    case task = "-."
}

class VxdayReader {
    
    static func allLists() -> [ListName] {
        let fm = FileManager.default
        let enumerator = fm.enumerator(atPath: VxdayFile.activeDir)!
        var lists: Set<String> = Set()
        while  let file  = enumerator.nextObject() as? String  {
            if file.characters.first != "." {
                lists.insert(VxdayUtil.beforeUnderscore(file)!)
            }
            
        }
        return lists.map {return ListName($0)}
    }
    static func hashToList(_ hash: Hash) -> ListName {
        return ListName("TODO")
    }
    
    static func itemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getSummaryFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        return items
    }
    static func completeItemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getCompleteFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        return items
    }
    static func readFile(_ path: String) -> [String] {
        guard let contents =  try? String(contentsOfFile: path) else {
            print("Error reading file: \(path)")
            return []
        }
        return contents.components(separatedBy: "\n").filter{ $0 != ""}
    }
    
    
    static func linesToItems(_ lines: [String], list: ListName) -> [Item] {
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
var now = VxdayUtil.now()

print("started on blah.")
now.addTimeInterval(-88000)
if let x = readLine(strippingNewline: true) {
    print("Read line \(x)")
}


let finish = VxdayUtil.now()
let interval = Int(finish.timeIntervalSince(now))
let hours = interval / 3600
let mins = (interval - hours * 3600) / 60
let seconds = (Int(interval)) % 60
print("thats \(hours), \(mins), \(seconds)")
print("done waiting. \(finish), thats \(interval)")


/*
 let location = "/Users/vic/Desktop/test.txt"
 let x =  try? String(contentsOfFile: location)
 print("x is \(x!)")

 */

/*
let allLists = VxdayReader.allLists()
print("all Lists: \(allLists)")
*/

//let list = ListName("vic")
//VxdayExec.allList(list)

/*
let summaryPath = VxdayFile.getSummaryFilename(list)
let contents = VxdayReader.readFile(summaryPath)
print("CONTENTS ARE: \(contents)'")
let items = VxdayReader.readSummary(contents, list: list)
print("items are: \(items)")
 */




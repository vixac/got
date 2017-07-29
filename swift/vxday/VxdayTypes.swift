//
//  VxdayTypes.swift
//  vxday
//
//  Created by vic on 28/07/2017.
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
    
    func english() -> String {
        
        switch self {
        case .job:
            return "Job"
        case .completeJob:
            return "Job Completed"
        case .task:
            return "Task"
        case .completeTask:
            return "Task Completed"
        case .token:
            return "Token"
        }
    }
}

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
    
    func pretty(dailyResolution: Bool = false ) -> String {
        return self.date.ago(dailyResolution: dailyResolution)
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
    
    func pretty(dailyResolution: Bool = false ) -> String {
        return self.date.ago(dailyResolution: dailyResolution)
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

struct  TimeBreakdown {
    let hours: Int
    let mins: Int
    let seconds : Int
    let totalSeconds: Int
    init(start: Date, end: Date) {
        
        totalSeconds = Int(end.timeIntervalSince(start))
        hours = totalSeconds / 3600
        mins = (totalSeconds - hours * 3600) / 60
        seconds = totalSeconds % 60
    }
    init(_ seconds: IntOffset) {
        totalSeconds = seconds.offset
        hours = totalSeconds / 3600
        mins = (totalSeconds - hours * 3600) / 60
        self.seconds = totalSeconds % 60
        
    }
    func toString() -> String {
        var str = ""
        if hours > 0 {
            str += "\(hours) hrs "
        }
        if mins > 0 {
            str += "\(mins) mins "
        }
        if seconds > 0 && !(hours > 0  && mins == 0) { //because 3 hrs 20 seconds is confusing on the eye face.
            str += "\(seconds) seconds"
        }
        return str
    }
}



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
        let times = TimeBreakdown(start: creation.date, end: completion.date)
        return "\(ItemType.token.rawValue) \(hash.hash) \(creation.toString()) \(times.hours) \(times.mins) \(times.seconds)"
    }
    
    func isComplete() -> Bool {
        return true
    }
    func itemType() -> ItemType {
        return .token
    }
    // TODO make lazy propery
    func timeBreakdown() -> TimeBreakdown {
        return TimeBreakdown(start: creation.date, end: completion.date)
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
            guard let hours = ArgParser.offset(args: array, index: 3) else {
                print("error getting hours from : \(array)")
                return nil
            }
            guard let mins = ArgParser.offset(args: array, index: 4) else {
                print("error getting minutes from : \(array)")
                return nil
            }
            guard let seconds = ArgParser.offset(args: array, index: 5) else {
                print("error getting seconds from : \(array)")
                return nil
            }
            var date = creationDate.date
            date = date.addingTimeInterval(TimeInterval(seconds.offset))
            date = date.addingTimeInterval(TimeInterval(mins.offset * 60))
            date = date.addingTimeInterval(TimeInterval(hours.offset * 3600))
            return Item.token(VxToken(list: list, hash: hash, creation: creationDate, completion: CompletionDate(date)))
            
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



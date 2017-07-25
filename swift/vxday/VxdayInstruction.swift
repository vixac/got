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

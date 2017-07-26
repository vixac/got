//
//  VxdayView.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation




enum Item {
    case job(ListName, Hash, CreationDate , DeadlineDate, Description )
    case completeJob(ListName, Hash, CreationDate, DeadlineDate, completion: CompletionDate)
    case task(ListName, Hash, CreationDate, Description)
    case token(ListName, Hash, CreationDate, CompletionDate)
    
    
    
    func toString() -> String {
        switch self {
            case let .job(_, hash, creation, deadline, description):
                return "\(ItemType.job.rawValue) \(hash.hash) \(creation.toString()) \(deadline.toString()) \(description.text)"
            case let .task(_, hash, creation, description):
                return "\(ItemType.task.rawValue) \(hash.hash) \(creation.toString()) \(description.text)"
            case let .token(_, hash, creation, completion):
                return "\(ItemType.tokenEntry.rawValue) \(hash.hash) \(creation.toString()) \(completion.toString())"
        default:
            return "TODO convert item to string."
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
            return Item.job(list, hash, creationDate, deadlineDate, description)
            
        case .task:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            
            guard let description = ArgParser.description(args: array, start: 3) else {
                print("Error: could not get description from: \(array)")
                return nil
            }
            return Item.task(list, hash, creationDate, description)
        case .tokenEntry:
            guard let creationDate = ArgParser.creation(args: array, index: 2) else {
                print("Error: could not extract creation date from: \(array)")
                return nil
            }
            guard let completionDate = ArgParser.completion(args: array, index: 2) else {
                print("Error: could not extract completion date from: \(array)")
                return nil
            }
            return Item.token(list, hash, creationDate, completionDate)
        default:
            print("TODO unhandled vxday line: \(line).")
            return nil
        }
    }
    
    
}


class VxdayView {
    
    
    
    
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




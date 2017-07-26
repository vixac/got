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
    static let baseColor : String = ANSIColors.reset.rawValue
    static let happyColor: String = ANSIColors.green.rawValue
    static let titleColor :String = ANSIColors.white.rawValue
    static let info2Color: String = ANSIColors.test.rawValue
    static let infoColor :String = ANSIColors.cyan.rawValue
    
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
        default:
             return "unknown"
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
        return VxdayColor.info(pad(hash.hash, toLength: 11))
    }
    
    func daysView(_ date: Date) -> String {
        let daysAgo = date.daysAgoInt()
        let agoStr = pad(date.daysAgo(), toLength: 14)
        let datedStr = agoStr + VxdayUtil.dateFormatter.string(from: date) + "   "
        if daysAgo < 0 {
            return VxdayColor.danger(datedStr)
        }
        if daysAgo == 0 {
            return  VxdayColor.warning(datedStr)
        }
        return  VxdayColor.happy(datedStr)
    }
    
    func listNameView(_ list: ListName) -> String {
        return VxdayColor.info2(pad(list.name, toLength: 9))
    }
    
    func showJobs() -> [String] {
        return self.getDeadlines().map {
            let daysAgo = $0.deadline.date.daysAgoInt()
            let timeStr =  pad( $0.deadline.pretty(), toLength: 15)
            var datedStr = timeStr +  VxdayUtil.dateFormatter.string(from: $0.deadline.date) + space()
            if daysAgo < 0 {
                datedStr = VxdayColor.danger(datedStr)
            }
            if daysAgo == 0 {
                datedStr = VxdayColor.warning(datedStr)
            }
            else {
                datedStr = VxdayColor.happy(datedStr)
            }
            
            let hash = hashView($0.hash)
            let listName = listNameView($0.list)
            return daysView($0.deadline.date) + hash + listName +  $0.description.text
            
        }
    }
    
    func oneLiners() -> [String] {
        if let listName = allLists().first?.name  {
            return listName
        }
        return "TODO ONE LINER"
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
        for _ in 1...needed {
            spaces = spaces +  " "
        }
        return string +  spaces
    }
}





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
        print("hashing this: \(hashable)")
        let hashed = "0" + String(format: "%2X", hashable).lowercased()
        print("hash is : \(hashed)")
        return hashed
    }
    
    class func increment(date: Date, byDays days : Int) -> Date {
        var comp = DateComponents()
        comp.day = days
        return  Calendar.current.date(byAdding: comp, to: date )!
    }
    
    
}


//TODO make this config
enum ItemType : String {
    case Complete = "x."
    case TokenEntry = "->."
    case Job = "=."
    case Task = "-."
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
    case yesterday = "yesterday"
    
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
    case week(IntOffset?)
    
    static func create(args: [String]) -> Instruction? {
        if args.count == 0 {
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
            print("TODOO FINISH UP")
            
            
            case .retire:
                guard let listName = ArgParser.listName(args: args, index: 1) else {
                    print("Error: Retire couldn't find list name in \(args)")
                    return nil
                }
                return .retire(listName)
            
            case .unretire:
                guard let listName = ArgParser.listName(args: args, index: 1) else {
                    print("Error: Unretire couldn't find list name in \(args)")
                    return nil
                }
                return .unretire(listName)
            
        }
    }
    
    
}

class ArgParser {
    static func listName(args: [String], index: Int) -> ListName? {
        if args.count < index {
            return nil
        }
        let str = args[index]
        if str.characters.first == "0" {
            print("Error, this looks like a hash, not a list name: \(str)")
            return nil
        }
        return ListName(str)
    }
    
    static func offset(args: [String], index: Int) -> IntOffset? {
        if args.count < index {
            return nil
        }
        let str : String = args[index]
        guard let intOffset : Int  = Int(str) else {
            print("Error, this doesn't look like an integer offset: \(str)")
            return nil
        }
        return IntOffset(intOffset)
    }
    
    static func description(args: [String], start: Int) -> Description? {
        guard let str = VxdayUtil.flattenRest(args, start: start) else {
            return nil
        }
        return Description(str)
    }
}

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

struct ListInstruction {
    let verb : Verb
    let name: ListName
    let offset: Int?
    let description: String?
    
    init?(verb: Verb, name: ListName, description: String? = nil, offset: Int? = nil ) {
        self.verb = verb
        self.name = name
        self.description = description
        self.offset = offset
        
        if verb == Verb.doIt  && description == nil {
            print("Error. Can't create a job without a description. Example: day do admin Send that email.")
            return nil
        }
        if verb == Verb.add && (description == nil || offset == nil) {
            print("Error. Can't add a job without a description and an offset in days. Example: day add home 7 Send gas meter reading.")
            return nil
        }
    }
    
    static func createJot(_ description: String)  -> ListInstruction {
        return ListInstruction(verb: .doIt, name: ListName("_jot"), description: description, offset: nil)!
    }
}

struct HashInstruction {
    let verb: Verb
    let name: ListName
    init(verb: Verb, name: ListName) {
        self.verb = verb
        self.name = name
    }
}

struct GlobalInstruction {
    let verb: Verb
    let offset:Int?
    init(verb: Verb, offset: Int?) {
        self.verb = verb
        self.offset = offset
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
        
        if type == .TokenEntry {
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

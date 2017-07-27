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

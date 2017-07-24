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
    
    static func hash(args: [String], index: Int) -> Hash? {
        if args.count < index {
            return nil
        }
        let str = args[index]
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
    
    
}

//
//  VxdayParser.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


struct IntOffset {
    let offset: Int
    init(_ offset: Int) {
        self.offset = offset
    }
}

enum Verb : String {
    case add = "till"
    case all = "jobs"
    case doIt = "to"
    case go = "go" //TODO RM
    case help = "help"
    case less = "less"
    case notes = "notes"
    case retire = "retire"
    case retired = "retired"
    case today = "today"
    case top = "top"
    case track = "track"
    case unretire = "unretire"
    case what = "what"
    case x = "done"
    case it = "it"
    case complete = "complete"
    case remove = "remove"
    case start = "start"
    case report = "report"
    case keep = "keep"
    case info = "info"
    
}

enum Instruction {
    
    //list actions
    case retire(ListName)
    case unretire(ListName)
    case add(ListName, IntOffset, Description)
    case doIt(ListName, Description)
    
    case lessList(ListName)
    case allList(String) //prefix, not list.
    case trackList(ListName)
    case top(ListName)
    case complete(ListName?)
    
    //hash actions
    case x(Hash)
    case go(Hash)
    case notes(Hash)
    case lessHash(Hash)
    case trackHash(Hash)
    case remove(Hash)
    case start(Hash)
    case info(Hash)
    //global actions

    case all
    case retired
    case what
    case track
    case week(IntOffset?)
    case report(IntOffset, ListName?)
    case help
    case keep(ListName)
    
    
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
            case .it:
                guard let description = ArgParser.description(args: args, start: 1) else {
                    print("Error: Do couldn't find a description in args: \(args)")
                    return nil
                }
                return .doIt(ListName("<nolist>"), description)
            case .add:
                var theOffset = IntOffset(0)
                if let offset = ArgParser.offset(args: args, index: 1) {
                    theOffset = offset
                } else if let dateOffset = ArgParser.dateOffset(args: args, index: 1) {
                    theOffset = dateOffset.offset
                } else {
                    guard let dateString = ArgParser.dateString(args: args, index: 1) else {
                        print("Error: Add couldn't find either a date offset or a readable date string. 4 means in 4 days, 4th means on the 4th of this or next month (whichever is in the futre)==")
                        return nil
                    }
                    theOffset =  dateString.offset
                }
                
                guard let listName = ArgParser.listName(args: args, index: 2) else {
                    print("Error: Add couldn't find list name in \(args)")
                    return nil
                }
                
                guard let description = ArgParser.description(args: args, start: 3) else {
                    print("Error: Add couldn't find a description in args: \(args)")
                    return nil
                }
                return .add(listName, theOffset, description)
                
            case .all:
                if let listName = ArgParser.listName(args: args, index: 1) {
                    return .allList(listName.name)
                }
                else {
                    return .all
                }
            case .info:
                guard let hash = ArgParser.hash(args: args, index: 1) else {
                    print("Error: Couldn't find hash  in \(args)")
                    return nil
                }
                return .info(hash)
            case .today:
                guard let listName = ArgParser.listName(args: args, index: 1) else {
                    print("Error: Do couldn't find list name in \(args)")
                    return nil
                }
                guard let description = ArgParser.description(args: args, start: 2) else {
                    print("Error: Do couldn't find a description in args: \(args)")
                    return nil
                }
                return .add(listName, IntOffset(0), description)
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
                
            case .notes:
                guard let hash = ArgParser.hash(args: args, index: 1) else {
                    print("Error: Couldn't find hash name in \(args)")
                    return nil
                }
                return .notes(hash)
        case .keep:
            guard let listName = ArgParser.listName(args: args, index: 1) else {
                print("Error: Do couldn't find list name in \(args)")
                return nil
            }
            return .keep(listName)
            case .start:
                guard let hash = ArgParser.hash(args: args, index: 1) else {
                    print("Error: Couldn't find hash name in \(args)")
                    return nil
                }
                return .start(hash)
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
            case .report:
                guard let offset = ArgParser.offset(args: args, index: 1) else {
                    print("Error: report couldn't extract number of days: day report 7 <list>")
                    return nil
                }
                return .report(offset, ArgParser.listName(args: args, index: 2))
            
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
  
    static func dateString(args: [String], index: Int) -> DateString? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return DateString(str)
    }
  
    static func dateOffset(args: [String], index: Int) -> DateOffset? {
        guard let str = ArgParser.str(args: args, index: index) else  {
            return nil
        }
        return DateOffset(str)
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
        if str.trimmingCharacters(in: CharacterSet.whitespacesAndNewlines)  == "" {
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

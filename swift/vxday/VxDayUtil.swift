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
        let x = String(format: "%2X", hashable).lowercased()
        print("hash is : \(x)")
        return x
        
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
    case x = "x"
    case retire = "retire"
    case unretire = "unretire"
    case less = "less"
    case track = "track"
    case all = "all"
    case today = "today"
    case go = "go"
    case note = "note"
    case yesterday = "yesterday"
    case what = "what"
    
}

struct ListInstruction {
    
}

struct HashInstruction {
    
}

struct GlobalInstruction {
    
}

struct Item {
    var type : ItemType
    var created: Date
    var completion: Date
    var description: String
    
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

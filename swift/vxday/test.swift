//
//  test.swift
//  vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//


import Foundation

class VxDayUtil {
    
    private static let datetimeFormat = "yyyy-MM-dd'T'HH:mm:ss"
    private static let dateFormat = "yyyy-MM-dd"
    
    static let datetimeFormatter : DateFormatter  = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = VxDayUtil.datetimeFormat
        return dateFormatter
    }()
    
    static let dateFormatter : DateFormatter = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = VxDayUtil.dateFormat
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
    
}


//TODO make this config
enum ItemType : String {
    case Complete = "x."
    case TokenEntry = "->."
    case Job = "=."
    case Task = "-."
}


struct Item {
    var type : ItemType
    var created: Date
    var completion: Date
    var description: String
    
    init?(_ line: String) {
        print("did we get here?.....")
        let array = VxDayUtil.splitString(line)
        guard array.count > 3  else {
            print("Error parsing. : \(line), not enough items in array: \(array)")
            return nil
        }
        guard let type = ItemType(rawValue: array[0]) else {
            print("Error parsing: \(line), unknown vxday type: \(array[0])")
            return nil
        }
        
        if type == .TokenEntry {
            print("TODO handle token entry")
            return nil
        }
        else {
            let hash = array[1]
            guard VxDayUtil.isValidHash(hash) else {
                print("invalid hash.: \(hash)")
                return nil
            }
            
            guard let createdTime = VxDayUtil.datetimeFormatter.date(from: array[2]) else {
                print("Error extracting created time: \(array[2])")
                return nil
            }
            guard let deadline = VxDayUtil.dateFormatter.date(from: array[3]) else {
                print("Error extracted completion date from : \(array[3])")
                return nil
            }
            guard let description = VxDayUtil.flattenRest(array, start: 4) else {
                print("Error, this job has no description.")
                return nil
            }
            self.type = type
            self.created = createdTime
            self.completion = deadline
            self.description = description
            print("did we get here?")
        }
    }
    func toString() -> String {
        let typeStr = type.rawValue
        let createdStr = VxDayUtil.datetimeFormatter.string(from: self.created)
        let completionStr = VxDayUtil.dateFormatter.string(from: self.completion)
        return typeStr + " " + createdStr + " " + completionStr + " " + description
    }
}


func printArgs() {
    for  argument in CommandLine.arguments {
        print(argument)
    }
}

//printArgs()
//print("Done Swifting.")
let original = "-. abcaa 2013-05-08T19:03:53 2013-05-08 One Two Three Four."
let item = Item(original)
print("original text is :")
print(original)
print("item is ")
print(item!)
print("converted back its:")
print(original)

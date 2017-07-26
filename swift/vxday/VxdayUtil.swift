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
    
    class func flattenRest(_ array: [String], start: Int) -> String? {
        if array.count < start {
            print("error flattening the rest.")
            return nil
        }
        let endWords = array.suffix(array.count - start)
        return endWords.flatMap({$0 + " " }).joined()
    }
    
    
    class func nowDay() -> Date {
        let c = Calendar.current
        return c.startOfDay(for: now())
    }
    
    class func now() -> Date {
        return Date()
    }
    
    class func hash(_ string: String) -> String {
        let time = now().timeIntervalSince1970.hashValue
        let hashable = "\(time)" + string
        var hashed = "0" + String(format: "%2X", hashable).lowercased()
        if hashed.characters.count < 9 {
            let numNeeded = 9 - hashed.characters.count
            for _ in 1...numNeeded {
                hashed += "0"
            }
        }
        return hashed
    }
    
    class func increment(date: Date, byDays days : Int) -> Date {
        var comp = DateComponents()
        comp.day = days
        return  Calendar.current.date(byAdding: comp, to: date )!
    }
    
    class func humanDuration(between start: Date, and end: Date) -> String {
        return "TODO"
        
    }
    
    class func beforeUnderscore(_ string: String) -> String {
        return string.components(separatedBy: "_").first ?? ""
    }
    
}


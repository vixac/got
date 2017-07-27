//
//  Vxday_util.swift
//  Vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class VxdayUtil {
    
    enum TimeBucket {
        case past
        case present
        case future
    }
    
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
    
    class func afterHyphen(_ string: String) -> String {
        return string.components(separatedBy: "/").last ?? string
    }
    
    class func beforeUnderscore(_ string: String) -> String? {
        return string.components(separatedBy: "_").first 
    }
    
}


extension Date {
    static func daysOffsetString(_ days: Int) -> String {
        switch days {
        case let x where x == -1:
            return "Yesterday"
        case let x where x == 0:
            return "Today"
        case let x where x == 1:
            return "Tomorrow"
        case let x where x < 0:
            return "\(abs(days)) days ago"
        default:
            return "In \(days) days"
        }
    }
    
    func daysAgoInt() -> Int {
        return Int(self.timeIntervalSince(VxdayUtil.nowDay())) / 86400
    }
    
    func daysAgo() -> String {
        return Date.daysOffsetString(self.daysAgoInt())
    }
    
    func bucket() -> VxdayUtil.TimeBucket {
        let daysAgo = self.daysAgoInt()
        if daysAgo < 0 {
            return .past
        }
        if daysAgo == 0 {
            return .present
        }
        return .future
    }
    
    func ago() -> String {
        let SECONDS_IN_A_DAY = 86400
        let SECONDS_IN_AN_HOUR = 60 * 60
        let SECONDS_IN_A_MINUTE = 60
        let now = VxdayUtil.now()
        
        
        let interval =  Int(self.timeIntervalSince(now))
        let absInterval = abs(interval)
        if absInterval > SECONDS_IN_A_DAY {
            return Date.daysOffsetString(interval / SECONDS_IN_A_DAY)
        }
        else if absInterval > SECONDS_IN_AN_HOUR {
            let hours = interval / SECONDS_IN_AN_HOUR
            if hours > 0 {
                return "In \(hours) hours."
            }
            
            return "\(abs(hours)) hours ago"
        }
        else if absInterval > SECONDS_IN_A_MINUTE {
            let mins = interval / SECONDS_IN_A_MINUTE
            if mins > 0 {
                return "In \(mins) mins."
            }
            return "\(abs(mins)) mins ago."
        }
        return "Just now."
    }
}


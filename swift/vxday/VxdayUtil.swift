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
    private static let timeFormat = "HH:mm:ss"
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
    
    static let timeFormatter: DateFormatter = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = VxdayUtil.timeFormat
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
    
    class func hash(_ string: String) -> Hash {
        let time = now().timeIntervalSince1970.hashValue
        let hashable = "\(time)" + string
        var hashed = "0" + String(format: "%2X", hashable).lowercased()
        if hashed.characters.count < 9 {
            let numNeeded = 9 - hashed.characters.count
            for _ in 1...numNeeded {
                hashed += "0"
            }
        }
        return Hash(hashed)
    }
    
    class func timelinessToString(_ daysOffset: Int) -> String {
        if daysOffset == -1 {
            return "1 day late"
        }
        if daysOffset < 0 {
            return "\(abs(daysOffset)) days late"
        }
        if daysOffset == 0 {
            return "On time"
        }
        if daysOffset == 1 {
            return "1 day early"
        }
        return "\(daysOffset) days early"
    }
    class func timeliness(deadline: DeadlineDate, completion: CompletionDate) -> Int {
        return  deadline.date.daysSince(completion.date)
        
    }
    class func increment(date: Date, byDays days : Int) -> Date {
        var comp = DateComponents()
        comp.day = days
        return  Calendar.current.date(byAdding: comp, to: date )!
    }

    
    class func afterHyphen(_ string: String) -> String {
        return string.components(separatedBy: "/").last ?? string
    }
    
    class func beforeUnderscore(_ string: String) -> String? {
        return string.components(separatedBy: "_").first 
    }
    class func pad(_ string: String, toLength length: Int) -> String {
        if string.characters.count > length {
            return string
        }
        var spaces = ""
        let needed = length - string.characters.count
        if needed < 1 {
            return string
        }
        for _ in 1...needed {
            spaces = spaces +  " "
        }
        return string +  spaces
    }
}

extension Date {
    static let SECONDS_IN_A_DAY = 86400
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
    func daysSince(_ date: Date) -> Int {
        return Int(self.timeIntervalSince(date)) / Date.SECONDS_IN_A_DAY
    }
    func daysAgoInt() -> Int {
        return self.daysSince(VxdayUtil.now())
    }
    
    func daysAgo() -> String {
        return Date.daysOffsetString(self.daysAgoInt())
    }
    
    func toDateString() -> String {
        return VxdayUtil.dateFormatter.string(from: self)
    }
    
    func toDateTimeString() -> String {
        return VxdayUtil.datetimeFormatter.string(from: self)
    }
    
    func toTimeString() -> String {
        return VxdayUtil.timeFormatter.string(from: self)
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
    
    func ago(dailyResolution: Bool = false) -> String {
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

    static public func from(year:Int, month:Int, day:Int) -> Date {
        var c = DateComponents()
        c.year = year
        c.month = month
        c.day = day
        
        let gregorian = Calendar(identifier:Calendar.Identifier.gregorian)
        let date = gregorian.date(from: c)
        return date!
    }
    
    public func isSameDayAs(_ date: Date) -> Bool {
        
        let dateComps = (Calendar.current as NSCalendar).components([.day, .month, .year], from: self)
        let otherDateComps = (Calendar.current as NSCalendar).components([.day, .month, .year], from: date)
        
        return dateComps.day == otherDateComps.day
            && dateComps.month == otherDateComps.month
            && dateComps.year  == otherDateComps.year
        
    }
    public func monthsToDate(_ date:Date) -> Int {
        let calendar = Calendar.current
        let components = (calendar as NSCalendar).components([.year, .month],
                                                             from: self,
                                                             to: date,
                                                             options: [])
        
        var months = 0
        months = components.year! * 12
        months = months + components.month!
        
        return components.month!
    }
    public func incrementByDays(_ numberOfDays:Int) -> Date {
        
        let calendar = Calendar(identifier: Calendar.Identifier.gregorian)
        var components = DateComponents()
        components.day = numberOfDays
        
        return (calendar as NSCalendar).date(byAdding: components,to:self, options: [])!
    }
    public func midnightsToDate(_ date:Date) -> Int {
        let calendar = Calendar.current
        let unitFlags : NSCalendar.Unit = [NSCalendar.Unit.year, NSCalendar.Unit.month, NSCalendar.Unit.day]
        let selfComponents = (calendar as NSCalendar).components(unitFlags, from: self)
        let selfMidnight = calendar.date(from: selfComponents)
        
        let dateComponents = (calendar as NSCalendar).components(unitFlags, from: date)
        let dateMidnight = calendar.date(from: dateComponents)
        
        let diffComponents = (calendar as NSCalendar).components( .day, from: selfMidnight!, to: dateMidnight!, options: [])
        return diffComponents.day!
    }
    public func daysToDate(_ date:Date) -> Int {
        let calendar = Calendar.current
        let unit:NSCalendar.Unit = .day
        let components = (calendar as NSCalendar).components(unit, from: self, to: date, options: [])
        return components.day!
    }
    
    public func maxoutDay() -> Date { //set the given time to 23:59:59
        let calendar = Calendar(identifier: Calendar.Identifier.gregorian)
        var components = (calendar as NSCalendar).components([NSCalendar.Unit.day, NSCalendar.Unit.month, NSCalendar.Unit.year],
                                                             from: self)
        
        components.hour = 23
        components.minute = 59
        components.second = 59
        
        return calendar.date(from: components)!
    }
    public func startOfDay() -> Date {
        let calendar = Calendar(identifier: Calendar.Identifier.gregorian)
        var components = (calendar as NSCalendar).components([NSCalendar.Unit.day, NSCalendar.Unit.month, NSCalendar.Unit.year],
                                                             from: self)
        
        components.hour = 0
        components.minute = 0
        components.second = 0
        
        return calendar.date(from: components)!
    }
    
    
}

extension String {
    public func wrap(columns: Int) -> [String] {
        let scanner = Scanner(string: self)
        
        var result : [String] = []
        var currentLineLength = 0
        
        var str = ""
        var word: NSString?
        while scanner.scanUpToCharacters(from: CharacterSet.whitespacesAndNewlines, into: &word) {
            let wordLength = word?.length ?? 0
            
            if currentLineLength != 0 && currentLineLength + wordLength + 1 > columns {
                // too long for current line, wrap
                str += " "
                result.append(str)
                str  = ""
                currentLineLength = 0
            }
            
            // append the word
            if currentLineLength != 0 {
                str += " "
                currentLineLength += 1
            }
            
            str += ((word ?? "") as String)
            currentLineLength += wordLength
        }
        result.append(str)
        return result
    }
}

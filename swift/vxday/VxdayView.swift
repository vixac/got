//
//  VxdayView.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class VxdayColor {
    static let dangerColor : String = ANSIColors.red.rawValue
    static let warningColor : String = ANSIColors.yellow.rawValue
    static let baseColor : String = ANSIColors.reset.rawValue
    static let happyColor: String = ANSIColors.green.rawValue
    static let titleColor :String = ANSIColors.white.rawValue
    static let info2Color: String = ANSIColors.test.rawValue
    static let infoColor: String = ANSIColors.cyan.rawValue
    static let whiteColor : String = ANSIColors.white.rawValue
    
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
    static func boldInfo(_ string: String) -> String {
        return whiteColor + string + baseColor
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
}


class Spaces {
    static let List = 9
    static let Timeliness = 14
    static let WhatOverdue = 13
    static let WhatPresent = 13
    static let WhatFuture = 13
    static let WhatTasks = 12
    static let DaysString = 15
    static let Hash = 11
    
}



class ListSummary {
    let list: ListName
    var past: [VxJob] = []
    var present: [VxJob] = []
    var future: [VxJob] = []
    var taskCount: Int = 0
    var timeWorked: Int = 0
    init(_ list: ListName) {
        self.list = list
    }
    func addJob(_ job: VxJob) {
        let bucket = job.deadline.date.bucket()
        switch bucket {
        case .past:
            past.append(job)
        case .present:
            present.append(job)
        case .future:
            future.append(job)
        }
    }
     func addTask(_ task: VxTask) {
        taskCount += 1
    }
    
    func addToken(_ token: VxToken) {
        let breakdown = token.timeBreakdown()
        timeWorked += breakdown.totalSeconds
    }
    func total() -> Int {
        return past.count + present.count + future.count + taskCount
    }
}


class DaySummary {
    let date: CreationDate
    var listSummaries: [ListName: IntOffset] = [:]
    var totalSeconds: Int = 0
    init(_ date: CreationDate ) {
        self.date = date
    }

    func addSomeSecondsTo(_ list: ListName, seconds: IntOffset) {
        totalSeconds += seconds.offset
        if let current = listSummaries[list] {
            listSummaries[list] = IntOffset(current.offset + seconds.offset)
        }
        else {
            listSummaries[list] = seconds
        }
    }
    func getSorted() -> [(list: ListName, duration: IntOffset)] {
        var tuples: [(list: ListName, duration: IntOffset)] = []
        for (list, duration) in listSummaries {
            tuples.append((list: list, duration:duration))
        }
        return tuples.sorted(by: {$0.duration.offset < $1.duration.offset})
    }
}

class TokenReport {
    
    //TODO make CreationDate hashable
    var days: [Date : DaySummary] = [:]
    
    func addToken(_ token: VxToken) {
        let seconds = IntOffset(token.timeBreakdown().totalSeconds)
        let day = token.creation.date.startOfDay()
        if let _ = days[day] {
            days[day]?.addSomeSecondsTo(token.list, seconds: seconds)
        }
        else {
            let summary = DaySummary(CreationDate(day))
            summary.addSomeSecondsTo(token.list, seconds: seconds)
            days[day] = summary
        }
    }
    
}
//TODO not used this yet:
class HashSummary {
    var hash: Hash
    var list: ListName
    var tokens: [VxToken] = []
    var totalSeconds : Int = 0
    init(_ list: ListName, hash: Hash) {
        self.list = list
        self.hash = hash
    }
    func addToken(_ token: VxToken) {
        guard hash.hash ==  token.hash.hash else {
            print("Dev error adding wrong token to Hash Summary. :\(hash.hash) isnt this tokens hash: \(token.hash.hash)")
            return
        }
        totalSeconds += token.timeBreakdown().totalSeconds
        tokens.append(token)
    }
}
class VxdayView {
    
    let items: [Item]
    init(_ item : Item) {
        items = [item]
    }
    init(_ items: [Item]) {
        self.items = items
    }
    
    func toBuckets() -> [ListSummary] {
        var dict : [ListName: ListSummary] = [:]
        for item in items {
            if let job = item.getJob() {
                let list = job.list
                if dict[list] == nil {
                     let summary = ListSummary(list)
                    summary.addJob(job)
                    dict[list] = summary
                }
                else  {
                    dict[list]?.addJob(job)
                }
                
            }
            else if let task = item.getTask() {
                let list = task.list
                if dict[list] == nil {
                    let summary =  ListSummary(list)
                    summary.addTask(task)
                    dict[list] = summary
                }
                else {
                    dict[list]?.addTask(task)
                }
                
            }
            else if let token = item.getToken() {
                let list = token.list
                if dict[list] == nil {
                    let summary =  ListSummary(list)
                    summary.addToken(token)
                    dict[list] = summary
                }
                else {
                    dict[list]?.addToken(token)
                }
                
            }
            //tokens aren't part of the summary yet.
        }
        return dict.map { return $0.value }.sorted( by: { $0.total() < $1.total()})
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

    func timeBucketToColor(_ date: Date, string: String)  -> String {
        let bucket = date.bucket()
        switch bucket {
            case .past:
                return VxdayColor.danger(string)
            case .present:
                return VxdayColor.warning(string)
            case .future:
                return  VxdayColor.happy(string)
            }
    }


    func showJobs() -> [String] {
        return self.getDeadlines().map {
            let timeStr =  pad( $0.deadline.pretty(), toLength: Spaces.DaysString)
            var datedStr = timeStr +  VxdayUtil.dateFormatter.string(from: $0.deadline.date) + spaces()
            datedStr = timeBucketToColor($0.deadline.date, string: datedStr)
            
            let hash = hashView($0.hash)
            let listName = listNameView($0.list)
            return daysView($0.deadline.date) + hash + listName +  $0.description.text
            
        }
    }

    func noStringForZero(_ prefix: String, number: Int, toLength : Int) -> String {
        if number == 0 {
            return pad("", toLength: toLength)
        }
        let str = "\(prefix) \(number)"
        return pad(str, toLength: toLength)
    }
    
    func globalOneLiner(buckets: [ListSummary]) -> String {
        let listCount = buckets.count
        let overdueCount = buckets.map { $0.past.count}.reduce(0, { $0 + $1})
        let todayCount = buckets.map { $0.present.count}.reduce(0, { $0 + $1})
        let futureCount = buckets.map { $0.future.count}.reduce(0 , {$0 + $1})
        let taskCount = buckets.map {$0.taskCount}.reduce(0, {$0 + $1})
        let totalLists = VxdayColor.boldInfo(pad("Lists: \(listCount)", toLength: Spaces.List))
        let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: overdueCount, toLength:  Spaces.WhatOverdue))
        let present = VxdayColor.warning(noStringForZero("Today:", number: todayCount, toLength: Spaces.WhatPresent))
        let future = VxdayColor.happy(noStringForZero("Upcoming:", number: futureCount, toLength: Spaces.WhatFuture))
        let tasks  = VxdayColor.warning(noStringForZero("Tasks:", number: taskCount, toLength: Spaces.WhatTasks))
        let total = VxdayColor.boldInfo("Total: \(overdueCount + todayCount + futureCount + taskCount)")
        return "\(totalLists) \(overdue) \(present) \(future) \(tasks) \(total)"
    }
    
    func oneLiners(_ buckets: [ListSummary]) -> [String] {
        return buckets.map  { summary in
            let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: summary.past.count, toLength: Spaces.WhatOverdue))
            let present = VxdayColor.warning(noStringForZero("Today:", number: summary.present.count, toLength: Spaces.WhatPresent))
            let upcoming = VxdayColor.happy(noStringForZero("Upcoming:", number: summary.future.count, toLength: Spaces.WhatFuture))
            let tasks = VxdayColor.warning(noStringForZero("Tasks:", number: summary.taskCount, toLength: Spaces.WhatTasks))
            let total = "Total: \(summary.past.count + summary.present.count + summary.future.count + summary.taskCount)"
            let listName = listNameView(summary.list)
            return "\(listName) \(overdue) \(present) \(upcoming) \(tasks) \(total)"
        }
    }
    
    func renderComplete() -> [String] {
        var strings: [String] = []
        items.forEach { item in
            var completed : CompletionDate? = nil
            var description : Description = Description("")
            var offsetFromDeadline : Int? = nil
            if let job = item.getJob() {
                completed = job.completion
                
                description = job.description
                if let c = completed {
                    offsetFromDeadline = VxdayUtil.timeliness(deadline: job.deadline, completion: c)
                }
            }
            else if let task = item.getTask() {
                completed = task.completion
                description = task.description
            }
            
            if let c = completed {
                //let overdue = VxdayColor.danger(noStringForZero("Overdue:", number: overdueCount, toLength:  Spaces.WhatOverdue))
                let completed = VxdayColor.boldInfo( pad(c.pretty(), toLength: Spaces.DaysString))
                let d = description.text
                
                var timelinessStr = ""
                if let o = offsetFromDeadline {
                    timelinessStr = VxdayUtil.timelinessToString(offsetFromDeadline!)
                    timelinessStr = pad(timelinessStr, toLength:  Spaces.Timeliness)
                    if o < 0 {
                        timelinessStr = VxdayColor.danger(timelinessStr)
                    }
                    else if o == 0 {
                        timelinessStr = VxdayColor.warning(timelinessStr)
                    }
                    else {
                        timelinessStr = VxdayColor.happy(timelinessStr)
                    }
                }
                else {
                    timelinessStr =  pad("", toLength: Spaces.Timeliness)
                }
                
                let hash = hashView(item.vxItem().hash)
                let list = listNameView(item.vxItem().list)
                strings.append("\(completed) \(timelinessStr) \(hash) \(list) \(d)")
            }
        }
        return strings
    }
    
    
    /*
    func renderItemStoredSummary() -> String {
        guard let item = items.first else {
            return ""
        }
        let vxItem = item.vxItem()
        let noun = vxItem.itemType().english()
        let hashStr =  hashView(vxItem.hash)
        let listStr = listNameView(vxItem.list)
        var due = ""
        if case ItemType.job = vxItem.itemType() {
     
            due = ", due \(daysView(item.getJob()!.deadline.date))"
        }
        return "\(noun) created with hash \(hashStr) for list \(listStr) \(due)"
    }
 */
    
    
    
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
  
    func spaces() -> String {
        return "   "
    }
    
    func daysView(_ date: Date) -> String {
        let agoStr = pad(date.daysAgo(), toLength: 14)
        let datedStr = agoStr + VxdayUtil.dateFormatter.string(from: date) + spaces()
        return timeBucketToColor(date, string: datedStr)
    }
    
    private  func hashView(_ hash: Hash) -> String {
        return VxdayColor.info2(pad(hash.hash, toLength: Spaces.Hash))
    }
    
    private  func listNameView(_ list: ListName) -> String {
        return VxdayColor.info(pad(list.name, toLength: Spaces.List))
    }
    
    private  func pad(_ string: String, toLength length: Int) -> String {
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





//
//  VxdayInstruction.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class VxdayInstruction {
    
    static func makeAddString(_ listName: ListName, description: Description, offset: IntOffset?) -> String {
        
        let now = VxdayUtil.now()
        let created = VxdayUtil.datetimeFormatter.string(from: now)
        
        let hash = VxdayUtil.hash(created + description.text)
        
        if let o = offset {
            
            let deadline = VxdayUtil.dateFormatter.string(from: VxdayUtil.increment(date: now, byDays: o.offset))
            
            return "\(ItemType.Job.rawValue) \(hash) \(created) \(deadline) \(description.text)"
        }
        else {
            return " \(ItemType.Task.rawValue) \(hash) \(description.text)"
        }
    }
    
}

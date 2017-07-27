//
//  VxdayInstruction.swift
//  vxday
//
//  Created by vic on 24/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

class VxdayInstruction {
    
    static func executeInstruction(_ instruction : Instruction) {
        switch instruction {
            case let .add(list, offset, description):
                let now = VxdayUtil.now()
                let created = CreationDate(now)
                let deadline = DeadlineDate(VxdayUtil.increment(date: now, byDays: offset.offset))
                let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
                let item = VxJob(list: list, hash: hash, creation: created, deadline: deadline, description: description , completion: nil)
                
                VxdayExec.storeItem(item)
            
            
            case let .doIt(list, description):
                let now = VxdayUtil.now()
                let hash = VxdayUtil.hash(VxdayUtil.datetimeFormatter.string(from: now) + description.text)
                let created = CreationDate(now)
                let item = VxTask(list: list, hash: hash, creation: created, description: description, completion: nil)
                VxdayExec.storeItem(item)
            case let .retire(list):
                VxdayExec.retire(list)
            case let .unretire(list):
                VxdayExec.unretire(list)
            case let .lessList(list):
                VxdayExec.lessList(list)
            case let .allList(list):
                VxdayExec.allList(list)
            case .all:
                VxdayExec.all()
            case .what:
                VxdayExec.what()
            case let .x(hash):
                VxdayExec.x(hash)
            
           // case let .x(hash):
           //     VxdayExec.x(hash)
            
            
        default:
             print("TODO handle instruction: \(instruction)")
        }
    }

    
    //TODO RM
    /*
    static func makeAddString(description: Description, offset: IntOffset?) -> String {
        
        let now = VxdayUtil.now()
        let created = VxdayUtil.datetimeFormatter.string(from: now)
        
        let hash = VxdayUtil.hash(created + description.text)
        
        if let o = offset {
            let deadline = VxdayUtil.dateFormatter.string(from: VxdayUtil.increment(date: now, byDays: o.offset))
            
            return "\(ItemType.job.rawValue) \(hash) \(created) \(deadline) \(description.text)"
        }
        else {
            return " \(ItemType.task.rawValue) \(hash) \(created) \(description.text)"
        }
    }
    */
}

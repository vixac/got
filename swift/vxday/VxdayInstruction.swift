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
                let string = VxdayInstruction.makeAddString(description: description, offset: offset)
                VxdayExec.append(list, content: string)
            
            case let .doIt(list, description):
                let string = VxdayInstruction.makeAddString(description: description , offset: nil)
                VxdayExec.append(list, content: string)
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
            
           // case let .x(hash):
           //     VxdayExec.x(hash)
            
            
        default:
             print("TODO handle instruction: \(instruction)")
        }
    }

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
    
}

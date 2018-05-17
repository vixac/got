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
                VxdayExec.createJob(list, offset: offset, description: description)
            case let .doIt(list, description):
                VxdayExec.createTask(list, description: description)
            case let .retire(list):
                VxdayExec.retire(list)
            case let .unretire(list):
                VxdayExec.unretire(list)
            case let .lessList(list):
                VxdayExec.lessList(list)
            case let .allList(prefix):
                VxdayExec.allList(prefix)
            case .all:
                VxdayExec.all()
        case let .complete(list):
            VxdayExec.showComplete(list)
            case .what:
                VxdayExec.what()
            case let .x(hash):
                VxdayExec.x(hash)
            case let .start(hash):
                VxdayExec.startTokenSession(hash)
        case let .remove(hash):
                VxdayExec.remove(hash)
        case let .report(days, list):
            VxdayExec.report(days, list: list )
        case let .notes(hash):
            VxdayExec.showNotes(hash)
        case let .info(hash):
            VxdayExec.info(hash)
        case .gotInfo:
            VxdayExec.gotInfo()
        case .help:
            VxdayExec.help()
        default:
             print("TODO handle instruction: \(instruction)")
        }
    }

}

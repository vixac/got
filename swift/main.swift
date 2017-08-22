#!/usr/bin/swift 
//
//  add.swift
//  vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

func main() {
    var args = CommandLine.arguments
    
    args.remove(at: 0) // thanks, but we dont need the filename.
    guard let instruction = Instruction.create(args) else {
        print("Invalid instruction. Try day help")
        return
    }
    VxdayInstruction.executeInstruction(instruction)
    VxColor.putBack()
}

main()



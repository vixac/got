//
//  add.swift
//  vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//





let args = CommandLine.arguments

if args.count < 4 {
    print("Error in day add. The format is day add <list> <days> <description. Example: To add a job to file your taxes next week under the list accounts, run this:  day add accounts 7 File taxes omg.")
}



//
//  test.swift
//  vxday
//
//  Created by vic on 23/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//


import Foundation

func printArgs() {
    for  argument in CommandLine.arguments {
        print(argument)
    }
}




/*
let original = "-. abcaa 2013-05-08T19:03:53 2013-05-08 One Two Three Four."
let item = Item(original)
print("original text is :")
print(original)
print("item is ")
print(item!)
print("converted back its:")
print(original)
printArgs()
*/
//shell("echo", "wehey this works")
//shell("ls", "-a")
//shell("touch", "testfile.txt")
//shell("./vic.sh", "from swift!")

let srcDir = VxdayExec.getEnvironmentVar("VXDAY_SRC_DIR")!
print("Src dir is \(srcDir)")
let command = srcDir + "/bash/retire.sh"


//VxdayExec.retire(ListName("me"))
VxdayExec.unretire(ListName("me"))
//shell("cat out.txt | xargs wc")

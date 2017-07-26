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
//VxdayExec.retire(ListName("me"))
//VxdayExec.unretire(ListName("me"))




/*
let str = VxdayInstruction.makeAddString(ListName("me"), description: Description("Check add string looks correct."), offset: IntOffset(4))
print("str is \(str)")
let task = VxdayInstruction.makeAddString(ListName("vxday"), description: Description("Get task strings wokring."), offset: nil)
print("task is \(task)")
VxdayExec.append(ListName("wehey"), content: str)
*/



//VxdayExec.note(ListName("bam"), hash: Hash("abcdefg"))



//print("waiting: \(now)")
//VxdayExec.wait(ListName("me"), hash: Hash("abcdefg"))
/*
if let x = readLine(strippingNewline: true) {
    print("Read line \(x)")
}


let finish = VxdayUtil.now()
print("done waiting. \(finish)")

 */

/*
 let location = "/Users/vic/Desktop/test.txt"
 let x =  try? String(contentsOfFile: location)
 print("x is \(x!)")

 */

/*
let allLists = VxdayReader.allLists()
print("all Lists: \(allLists)")
*/

let list = ListName("vic")
VxdayExec.allList(list)

/*
let summaryPath = VxdayFile.getSummaryFilename(list)
let contents = VxdayReader.readFile(summaryPath)
print("CONTENTS ARE: \(contents)'")
let items = VxdayReader.readSummary(contents, list: list)
print("items are: \(items)")
 */




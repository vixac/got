//
//  VxdayRead.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

//TODO make this config
enum ItemType : String {
    case complete = "x."
    case completeJob = "X."
    case tokenEntry = "->."
    case job = "=."
    case task = "-."
}

class VxdayReader {
    
    static func allLists() -> [ListName] {
        let fm = FileManager.default
        let enumerator = fm.enumerator(atPath: VxdayFile.activeDir)!
        var lists: Set<String> = Set()
        while  let file  = enumerator.nextObject() as? String  {
            if file.characters.first != "." {
                lists.insert(VxdayUtil.beforeUnderscore(file)!)
            }
            
        }
        return lists.map {return ListName($0)}
    }
    static func hashToList(_ hash: Hash) -> ListName {
        return ListName("TODO")
    }
    
    static func itemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getSummaryFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.readSummary(contents, list: list)
        return items
    }
    
    static func readFile(_ path: String) -> [String] {
        guard let contents =  try? String(contentsOfFile: path) else {
            print("Error reading file: \(path)")
            return []
        }
        return contents.components(separatedBy: "\n").filter{ $0 != ""}
    }
    
    
    static func readSummary(_ lines: [String], list: ListName) -> [Item] {
        return lines.flatMap{ Item.create($0, list: list)}
    }

}

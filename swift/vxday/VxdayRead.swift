//
//  VxdayRead.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


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
    
    
    
    static func itemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getSummaryFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        return items
    }
    static func completeItemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getCompleteFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        return items
    }
    static func tokensForList(_ list: ListName) -> [VxToken] {
        let filename = VxdayFile.getTokenFilename(list)
        let contents = VxdayReader.readFile(filename)
        //danger we might get silent errors here with the flat map if getToken fails.
        return VxdayReader.linesToItems(contents, list: list).flatMap { $0.getToken()}
    }
    static func readFile(_ path: String) -> [String] {
        guard let contents =  try? String(contentsOfFile: path) else {
            print("Error reading file: \(path)")
            return []
        }
        return contents.components(separatedBy: "\n").filter{ $0 != ""}
    }
    
    
    static func linesToItems(_ lines: [String], list: ListName) -> [Item] {
        return lines.flatMap{ Item.create($0, list: list)}
    }

}

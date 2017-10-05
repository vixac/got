//
//  VxdayRead.swift
//  vxday
//
//  Created by vic on 26/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


class Cache {
    
    var descriptions: [Hash: Description] = [:]
    
    func addItems(_ items: [Item]) {
        items.forEach {
            if let j = $0.getJob() {
                descriptions[j.hash] = j.description
            }
            else if let t = $0.getTask() {
                descriptions[t.hash] = t.description
            }
        }
    }
}

class VxdayReader {
    
    //ugly.
    static var cache : Cache = Cache()
 
    static func allLists(_ prefix: String = "") -> [ListName] {
        let fm = FileManager.default
	print("VX looking for active:")
        guard let enumerator = fm.enumerator(atPath: VxdayFile.activeDir) else { return [] 
}
        var lists: Set<String> = Set()
        while  let file  = enumerator.nextObject() as? String  {
  	    if prefix != "" {
		if !file.hasPrefix(prefix) {
		     continue
                }
            }
            if file.characters.first != "." {
                lists.insert(VxdayUtil.beforeUnderscore(file)!)
            }
            
        }
        return lists.map {return ListName($0)}
    }

    static func isListPresent(_ list: ListName) -> Bool {
        return allLists().filter{ $0.name == list.name}.count > 0
    }
    static func itemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getSummaryFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        cache.addItems(items)
        return items
    }
    static func completeItemsInList(_ list: ListName) -> [Item] {
        let filename = VxdayFile.getCompleteFilename(list)
        let contents = VxdayReader.readFile(filename)
        let items = VxdayReader.linesToItems(contents, list: list)
        cache.addItems(items)
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
            //print("Error reading file: \(path)")
            return []
        }
        return contents.components(separatedBy: "\n").filter{ $0 != ""}
    }
    
    static func allTokensAfter(_ date: Date, lists: [ListName]) -> [VxToken] {
        var tokens: [VxToken] = []
        let intervalStart = date.timeIntervalSinceNow
        lists.forEach { list in
            VxdayReader.tokensForList(list)
                .filter { $0.creation.date.timeIntervalSince1970 >  intervalStart}
                .forEach {
                    tokens.append($0)
            }
        }
        return tokens
    }
    
    static func descriptionForHash(_ hash: Hash, list: ListName) -> Description? {
        if let d = cache.descriptions[hash] {
            return d
        }
        else {
            //cache up active tokens
            let _ = self.itemsInList(list)
        }
        if let d = cache.descriptions[hash] {
            return d
        }
        else {
            //cache up completed tokens (this is slower so we procrastinate doing this)
            let _ = self.completeItemsInList(list)
        }
        guard  let d = cache.descriptions[hash] else {
            print("Error, looked in summary and complete but couldnt find any description for hash: \(hash.hash)")
            return nil
        }
        return d
    }
    
    static func linesToItems(_ lines: [String], list: ListName) -> [Item] {
        return lines.flatMap{ Item.create($0, list: list)}
    }
}

//
//  main.swift
//  vxday
//
//  Created by vic on 19/03/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation

struct DayItem {
    var content: String
    var list: List
    var day: Day
}
struct List {
    var name: String
}

struct Day {
    var offset : Int
}

/*
  The file store. load and save files and thats it.
 */
protocol DayStore {
    func getItem( day: Day, list: List?)
    func saveItem(item: DayItem)
}

struct Prefix {
    var value: String
}

/*
  The query. Returns quereis.
 */
protocol DayQuery {
    func getItems(list: List?, prefix: Prefix)
}

enum DQLQuery {
    case getItem( day: Day, list: List?)
    case getAll(list: List?, prefix: Prefix)
}

enum QueryResult {
    
    case invalidQuery
}

protocol DQLParser {
    func argsToQuery(_ args: [String]) -> QueryResult
}


print("Hello Gorgeous!!!")






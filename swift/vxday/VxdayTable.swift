//
//  VxTable.swift
//  vxday
//
//  Created by vic on 29/07/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import Foundation


//TODO RM. Going to do this dynamically.
class Spaces {
    static let List = 9
    static let Timeliness = 14
    static let WhatOverdue = 13
    static let WhatPresent = 13
    static let WhatFuture = 13
    static let WhatTasks = 12
    static let DaysString = 15
    static let DaysAgo = 14
    static let Hash = 11
}


enum Cell {
    case list(ListName?)
    case deadline(DeadlineDate?) // yesterday 10 days ago in 10 days today etc.
    case yyyymmdd(Date?, VxColor)
    case hash(Hash?)
    case description(Description?)
    case created(CreationDate?)
    case timeliness(TimeBreakdown?)
    case dayCount(IntOffset?) //Today: 1, Overdue: 15 etc
    case overdue(Int?)
    case today(Int?)
    case upcoming(Int?)
    case total(Int?)
    case text(String?, VxColor)
    case totalTime(TimeBreakdown?) //1 hour, 20mins, etc
    
    func color() -> VxColor {
        switch self {
            case .list:
                return VxColor.info()
            case let .deadline(date):
                guard let d = date else {
                    return VxColor.base()
                }
                return Cell.timeBucketToColor(d.date.bucket())
            case let .yyyymmdd(_, color):
                return color
            case let .created(date):
                guard let d = date else {
                    return VxColor.base()
                }
                return Cell.timeBucketToColor(d.date.bucket())
        case .overdue:
            return VxColor.danger()
        case .today:
            return VxColor.warning()
        case .upcoming:
            return VxColor.happy()
        case .hash:
            return VxColor.info2()
        case .total:
            return VxColor.boldInfo()
        case let .text(_, color):
            return color
        case .timeliness:
            return VxColor.info2()
        
        default:
            return VxColor.base()
        }
    }
    
    func plainText() -> String {
        switch self {
        case let .list(l):
            return l?.name ?? ""
        case let .deadline(d):
            return d?.pretty() ?? ""
        case let .yyyymmdd(d, _):
            guard let date = d else {
                return ""
            }
            return VxdayUtil.dateFormatter.string(from: date)
        case let .hash(h):
            return h?.hash ?? ""
        case let .description(d):
            return d?.text ?? ""
        case let .created(c):
            return c?.pretty() ?? ""
        case let .timeliness(b):
            guard let breakdown = b else {
                return ""
            }
            return breakdown.toString()
        case let .overdue(v):
            return emptyOrPrefix("Overdue", v)
        case let .today(v):
            return emptyOrPrefix("Today", v)
        case let .upcoming(v):
            return emptyOrPrefix("Upcoming", v)
        case let .total(v):
            return emptyOrPrefix("Total", v)
        case let .text(text, _):
            return text ?? ""
        default:
            return "TODO Text: \(self)"
        }
    }
    
    private func emptyOrPrefix(_ prefix: String, _ val: Int?) -> String {
        return val == nil ? "" : (val! == 0 ? "" : ("\(prefix): \(val!)"))
    }
    private static func timeBucketToColor(_ bucket: VxdayUtil.TimeBucket) -> VxColor {
        switch bucket {
            case .past:
                return VxColor.danger()
            case .present:
                return VxColor.warning()
            case .future:
                return  VxColor.happy()
        }
    }
}


struct ReadyCell {
    let cell: Cell
    let plainText: String
    let color: VxColor
    let width: Int
    init(_ cell: Cell) {
        self.cell = cell
        self.plainText = cell.plainText()
        self.width = self.plainText.characters.count
        self.color = cell.color()
        
        
    }
}


enum Row {
    case cells([ReadyCell])
    case heading(title: String, char: String, color: VxColor)
}


class VxdayTable {
    //takes semantic information about each row.
    //wait. it doesnt have to. just a STRING i guess and a maybe a callback for the color.

    var columnWidths: [Int:Int] = [:] //column number to width.
    
    var rows: [Row] = []
    
    var title: String
    var columnNames: [String] = []
    init(title: String) {
        self.title = title
    }
    func addColumnTitles(_ titles: [String]) {
        columnNames = titles
    }
    func renderSectionDivider(_ text: String, char: String, totalLength: Int, color: VxColor) ->  String {
        
        let textLen = text.characters.count
        if totalLength <  textLen {
            return text
        }
        let repeatingSymbol = char
        let firstBitLength = (totalLength - textLen) / 2
        let firstBit = String(repeating: repeatingSymbol, count: firstBitLength / char.characters.count )
        var str = firstBit + text
        
        
        //padding here to keep the repeating pattern in sync.
        let repeatingLength = char.characters.count
        while str.characters.count % repeatingLength != 0 {
            str += " "
        }
        
        
        let remainingLength = totalLength - str.characters.count
        str += String(repeating: repeatingSymbol, count: remainingLength / char.characters.count )
        return color.colorThis(str)
    }
    
    func addHeading(_ title: String, char: String, color: VxColor) {
        rows.append(Row.heading(title: title, char: char, color: color))
    }
    func addRow(_ cells: [Cell]) {
        
        var readies : [ReadyCell] = []
        for (i , cell) in cells.enumerated() {
            let ready = ReadyCell(cell)
            self.newCellOnColumn(ready, column: i)
            readies.append(ready)
        }
        let row = Row.cells(readies)
        rows.append(row)
    }
    
    func newCellOnColumn(_ cell: ReadyCell, column: Int) {
        let spaceBetweenCells = 3
        let width  = cell.width + spaceBetweenCells
        if columnWidths[column] == nil {
            columnWidths[column] = width
        }
        if let w = columnWidths[column], w < width {
            columnWidths[column] = width
        }
    }
    
    
    private func renderTitle(length: Int) -> String {
        return self.renderSectionDivider(self.title, char: "=", totalLength: length, color: VxColor.boldInfo())
    }
    
    private func renderColumnNames() -> String {
        var str = ""
        for (i, name) in columnNames.enumerated() {
            let width = columnWidths[i]!
            str += VxdayTable.renderText(name, width: width, color: VxColor.boldInfo())
        }
        return str
    }
    
    func render() -> [String] {
        let tableWidth = columnWidths.map {$0.value}.reduce(0, {$0 + $1})
        
        var rendered: [String] = []
        if self.title != "" {
            rendered.append(renderTitle(length: tableWidth))
        }
        rendered.append(renderColumnNames())
        rows.forEach { row in
            
            var rowText = ""
            switch row {
                case let .cells(cells):
                    for (i, cell) in cells.enumerated() {
                        rowText += VxdayTable.renderText(cell.plainText, width: columnWidths[i]!, color: cell.color)
                }
                case let .heading(title, char, color):
                    rowText =  renderSectionDivider(title, char: char, totalLength: tableWidth, color: color)
                
            }

            rendered.append(rowText)
        }
        return rendered
    }
    
    static func renderText(_ text: String, width: Int, color: VxColor) -> String {
        let text = VxdayUtil.pad(text, toLength: width)
        return color.colorThis(text)
    }
    
}

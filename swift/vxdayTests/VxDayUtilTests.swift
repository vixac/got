//
//  vxdayTests.swift
//  vxdayTests
//
//  Created by vic on 24/03/2017.
//  Copyright © 2017 vixac. All rights reserved.
//

import XCTest

@testable import vxday

class VxDayUtilTest: XCTestCase {
    
    override func setUp() {
        super.setUp()
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }
    
    func testExample() {
        
        XCTAssertEqual(1, 20)
        // This is an example of a functional test case.
        // Use XCTAssert and related functions to verify your tests produce the correct results.
    }
    
    func testSplitString() {
        
        let array = VxDayUtil.splitString(string: "one two three")
        XCTAssertEqual(array.count, 3)
        
    }
}

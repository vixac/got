// swift-tools-version:5.1
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "Got",
    products: [
       .executable(name: "Got", targets: ["Got"])
    ],
    dependencies: [],
    targets: [
        .target(name: "Got", dependencies: [], path: "Sources"),
        .testTarget(name: "GotTests", dependencies: ["Got"], path: "Tests")
    ]
)


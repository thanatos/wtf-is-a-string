import Foundation

let objCString:NSString = NSString(characters: [0xd83d], length: 1)
print("Debugging the original string:")
for index in 0..<objCString.length {
    let c = objCString.character(at: index)
    print(String(format: "0x%x", c), " ", terminator: "")
}
print()
print("Done w/ NSString")

let swiftStr = "|" + (objCString as String) + "|"
print("swiftStr = ", swiftStr)
print("length = ", swiftStr.utf16.count)  // Emits 3
print("code units = [", terminator: "")
let first = true;
for codeUnit in swiftStr.utf16 {
    if first {
        first = false
    } else {
        print(", ", terminator: "")
    }
    print(String(format: "0x%x", codeUnit), terminator: "")
}
print("]")

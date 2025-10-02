// Test JSON.stringify and JSON.parse

print("=== Testing JSON Builtin ===\n")

// Test 1: stringify an object
print("1. JSON.stringify(object)")
let obj = { name: "Alice", age: 30, active: true }
let jsonStr = JSON.stringify(obj)
print("Original:", obj)
print("Stringified:", jsonStr)
print("")

// Test 2: stringify an array
print("2. JSON.stringify(array)")  
let arr = { a: 1, b: 2, c: 3 }
let arrStr = JSON.stringify(arr)
print("Array:", arr)
print("Stringified:", arrStr)
print("")

// Test 3: parse a JSON string
print("3. JSON.parse(string)")
let jsonString = '{"name":"Bob","age":25,"active":false}'
let parsed = JSON.parse(jsonString)
print("JSON string:", jsonString)
print("Parsed:", parsed)
print("Name:", parsed.name)
print("Age:", parsed.age)
print("Active:", parsed.active)
print("")

// Test 4: round-trip conversion
print("4. Round-trip conversion")
let original = { city: "New York", population: 8000000 }
let stringified = JSON.stringify(original)
let roundtrip = JSON.parse(stringified)
print("Original:", original)
print("After round-trip:", roundtrip)
print("City:", roundtrip.city)
print("Population:", roundtrip.population)
print("")

// Test 5: stringify primitives
print("5. Stringify primitives")
print("String:", JSON.stringify("hello"))
print("Number:", JSON.stringify(42))
print("Boolean:", JSON.stringify(true))
print("Null:", JSON.stringify(nil))
print("")

print("✅ All JSON tests completed!")

// GET request to list posts
print("=== GET List of Posts ===")
let response2 = fetch("https://jsonplaceholder.typicode.com/posts")
print("Status:", response2.status)
print("OK:", response2.ok)
print("Body length:", response2.body)
print("")

print("=== POST Request with JSON Body ===")
let response3 = fetch("https://jsonplaceholder.typicode.com/posts", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    body: '{"title": "My Post", "body": "This is a test post", "userId": 1}'
})
print("Status:", response3.status)
print("OK:", response3.ok)
print("Body:", response3.body)
print("")

print("=== Request with Custom Headers ===")
let response6 = fetch("https://jsonplaceholder.typicode.com/posts/1", {
    headers: {
        "Accept": "application/json",
        "User-Agent": "GoScript/1.0"
    }
})
print("Status:", response6.status)
print("OK:", response6.ok)
print("Response Headers:", response6.headers)
print("")

print("=== Error Handling - Invalid URL ===")
let response7 = fetch("not-a-valid-url")
if (response7.error) {
    print("Error:", response7.error)
}
print("")

print("All fetch tests completed!")

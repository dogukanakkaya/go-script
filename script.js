print("--- BASICS ---");

var arr = [10, 20, 30, 40, 50];
print("Array:", arr);
print("arr[0]:", arr[0]);
print("arr[2]:", arr[2]);
print("arr[4]:", arr[4]);
print("arr.length:", arr.length);
print("Out of bounds arr[10]:", arr[10]);
print("Negative index arr[-1]:", arr[-1]);

var nested = [[1, 2], [3, 4]];
print("Nested array:", nested);
print("nested[0]:", nested[0]);
print("nested[0][1]:", nested[0][1]);

let obj = { a: 1, b: 2, c: 3 };
print(obj);

var x = 10;
var y = 20;
var sum = x + y;
print("x =", x, ", y =", y, ", x + y =", sum);

var greeting = "Hello";
var name = "World";
var message = greeting + ", " + name + "!";
print(message);

var multiply = function(a, b) {
    return a * b;
};
var result = multiply(6, 7);
print("6 * 7 =", result);

var age = 25;
if (age >= 18) {
    print("You are an adult");
} else {
    print("You are a minor");
}

print("Counting from 1 to 5:");
var i = 1;
while (i <= 5) {
    print(i);
    i = i + 1;
}

var person = {
    name: "Alice",
    age: 30,
    city: "New York"
};
print("Name:", person.name);
print("Age:", person.age);
print("City:", person.city);

var factorial = function(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
};
print("factorial(5) =", factorial(5));

var makeCounter = function() {
    var count = 0;
    return function() {
        count = count + 1;
        return count;
    };
};
var counter = makeCounter();
print("counter() =", counter());
print("counter() =", counter());

print("true:", true);
print("false:", false);
print("!true:", !true);
print("!false:", !false);

print("5 == 5:", 5 == 5);
print("5 != 3:", 5 != 3);
print("10 > 5:", 10 > 5);
print("3 < 7:", 3 < 7);

print("--- /BASICS ---");


print("--- fetch/JSON ---");

print("=== GET List of Posts ===")
let response = fetch("https://jsonplaceholder.typicode.com/posts")
print("Status:", response.status)
print("OK:", response.ok)
print("Body length:", response.body)
print("")

print("=== POST Request with JSON Body ===")
let response2 = fetch("https://jsonplaceholder.typicode.com/posts", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    body: JSON.stringify({"title": "My Post", "body": "This is a test post", "userId": 1})
})
print("Status:", response2.status)
print("OK:", response2.ok)
print("Body:", response2.body)
print("")

print("=== Request with Custom Headers ===")
let response3 = fetch("https://jsonplaceholder.typicode.com/posts/1", {
    headers: {
        "Accept": "application/json",
        "User-Agent": "GoScript/1.0"
    }
})
print("Status:", response3.status)
print("OK:", response3.ok)
print("Response Headers:", response3.headers)
print("")

print("=== Error Handling - Invalid URL ===")
let response4 = fetch("not-a-valid-url")
if (response4.error) {
    print("Error:", response4.error)
}
print("")

print("--- /fetch/JSON ---");

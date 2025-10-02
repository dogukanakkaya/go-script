print("1. Variables and Arithmetic:");
var x = 10;
var y = 20;
var sum = x + y;
print("x =", x, ", y =", y, ", x + y =", sum);

print("2. String Operations:");
var greeting = "Hello";
var name = "World";
var message = greeting + ", " + name + "!";
print(message);

print("3. Functions:");
var multiply = function(a, b) {
    return a * b;
};
var result = multiply(6, 7);
print("6 * 7 =", result);

print("4. Conditionals:");
var age = 25;
if (age >= 18) {
    print("You are an adult");
} else {
    print("You are a minor");
}

print("5. Loops:");
print("Counting from 1 to 5:");
var i = 1;
while (i <= 5) {
    print(i);
    i = i + 1;
}

print("6. Objects:");
var person = {
    name: "Alice",
    age: 30,
    city: "New York"
};
print("Name:", person.name);
print("Age:", person.age);
print("City:", person.city);

print("7. Recursion - Factorial:");
var factorial = function(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
};
print("factorial(5) =", factorial(5));

print("8. Closures:");
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
print("counter() =", counter());

print("9. Boolean Operations:");
print("true:", true);
print("false:", false);
print("!true:", !true);
print("!false:", !false);

print("10. Comparisons:");
print("5 == 5:", 5 == 5);
print("5 != 3:", 5 != 3);
print("10 > 5:", 10 > 5);
print("3 < 7:", 3 < 7);

print("Completed");

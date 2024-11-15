## 030 Simple comparisons
Note: While in many languages, if two objects that have the same fields but have different addresses, they're not equal if they're
compared using == operator. That's because in those langs they're trying to differentiate between checking if memory addresses are
equal(using ==) vs checking if underlying values are equal.

But in go that's not the case. Whenever you have two structs, if they have same fields and values, they are equal, even though
they're memory addresses are different.

**Note: Structs containing functions as their fields can't be compared.** So you can't compare structs like these:
```go
package main

type Dog struct {
	// ...
	bark func() // this struct can't be compared with other structs because of this field
}
```

Note: You can get the memory location of a var using fmt.Print(**%p**, &x).

Pointers are also comparable. If they point to the same variable(if they point to the same memory location), they're equal.
Also, if both are nil, they're equal. Note that go doesn't look at the memory location of the pointer variable. It looks at where it **points**.

## 031 Reflects DeepEqual function
Look at 30 folder.

Look at `DogWithFn` type. It has a func field, so we can't use it with == operator, it will give compile err.

Note: There are cases where using reflect.DeepEqual() won't work.

## 032 Golden files brief overview
Another common technique for comparing test results is to use golden files.

### Golden files
TL;DR - covered in a later section, but here's a rough idea of what they are ...

If we need to compare big files(large csv, image, etc) then trying to recreate the `want` variable can be hard.
For example when altering an img, a lot of times we don't care about the underlying bytes as the result. Instead,
let's say for example we're making an img black and white, the easiest way to verify the code worked correctly, is to just look at
the resulting img. But we still need to somehow test it but computers can't easily look at it and verify it.

A common solution is to store a "golden file" - a file representing the desired test output - in our actual test source control
and to just compare to it directly. So we won't store that megabytes of data into our actual go code. Instead, we save and commit
a file containing those data.

For this, apply the code to an img, visually verify the result on that img, save that resulting img in the source control and commit it.
Then in future, we would have our tests take the same initial img, run the code(apply the filter) and make sure we still get the
same resulting img with what we stored in the source code.

Note: If we change the source code, the test is gonna fail but this solution is still nice because anytime we do make a change,
it forces the dev to visually verify the result is still correct, while still giving us a way to verify that things are still
working if we were not expecting to change that test.

When we learn about golden files(later section) we'll see tips/tricks to help manage this.

## 033 Helper comparison functions
### Helper funcs for comparisons

Sometimes in tests comparing the entire objects is not feasible. So we don't wanna recreate the entire res. We just care about specific attrs.
For example in http reqs, we only check for status codes. Or specific fields in the json res.

Solution: Write helper funcs.

In this pattern, we write `checkFn`s to pass a slice of these into testcases to check for various things. With this pattern,
makes it easier to check for specific things that we care about while still writing the tests that still use the same shared setup teardown.

Note: If the test helper funcs are becoming complex, you can take the test helpers(like `getText` and `findNodes`) themselves and
move them to their own packages. And then we can write tests for that package.
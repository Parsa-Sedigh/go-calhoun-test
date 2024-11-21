# Section 13 External and internal testing

## 048 Differences between external and internal tests
- Internal: you're inside the same package you wanna test. So you have access to everything inside that package. 
- External: We're not inside a package we wanna test. So we only have access to exported things.

This is the same for source code and tests.

But the difference is in tests, we have the freedom of choosing if we want to be external or internal.

## 049 How to write internal and external tests
When you have a package that has both internal and external tests, a common naming scheme is to use `_internal_test.go` as
the suffix for internal tests and for external test, don't use anything other than `_test.go`. The reason for this is
most people tend to focus on external tests, so there only would be a couple of internal tests, so it's nice to make it clear
that those are the internal tests.

Note: Let's say we're in package x. Internal tests use the same package as the source code they're testing, so the package for
those files is `x`, but the package name of external tests is `x_test`. By using x_test as package name of external tests,
it signifies to other devs that the only reason this is in another package, is so that we can test from another package(hence
being an external test).

Note: Both of the internal and external tests are in the same directory, but their packages are different. External tests package is: `x_test`.

So since the external tests are not in the same package, like other code that's not in the same package, they have to use
prefix of the package they are importing.

Do not import the packages using `.`. Like using it for importing the package that is defined in the same directory as the external
test. It makes it confusing. Because we made it an external test by using x_test as the package name, but now you're treating
it being an internal test.

Note: One of the use cases of tests is as a way for other devs to look at how to use your code as an example.

Note: When you have source code in the same directory, they all have to be in the same package. So if you have two files
in the same directory, they both have to be in the same package. If one of the files have a different package, it would get
compiler error(go build will throw an error) saying: `can't load packages: package x: found packages x and y in ...`.

But if you run `go test`, it would compile! The reason that works is testing tool has an exception to the mentioned rule where
if you have two different packages in your testing files, **that's acceptable**, because they wanted to make it easier for 
you to test from an external view point without needing to put that code in another directory because if you put the code in
another directory, all of the sudden it's not clear what code that test is even for, but by keeping it in the same directory,
it's clear: ok this is the directory that it's testing for.

**So we can actually have two `_test.go` files in the same directory with two different packages! But source files in the same directory
must have the same package.**

## 050 When to use external tests
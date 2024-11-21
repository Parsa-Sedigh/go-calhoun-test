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
**By default, write external tests.** Whenever you start to write a new test, add into `x_test` package(**in the same directory as the x package though**).
If you ever get to a position where you can't test sth, or change some data that you need to, because it's private, you still can
export those things you need to, so that you can access them in the external test but **only in the external test** and not in
the other source code files and then in the next vids we're also gonna look at why you might **still** need to write internal tests
and some specific reasons why they're useful and necessary at times.

**Anytime you can use an external test, you should.** So always go with external tests unless you can't.

That's counterintuitive because up to this point, we've always written internal tests.

External tests are gonna give you better tests. By writing an external test, your test is much less likely to break and it's much
more likely to last longer.

Why?

1. forces you to write test that is mostly agnostic as to how the actual work is getting done. Because we're not inside of the package
we're testing, so we don't have access to private things, therefore we're not aware of the details.
2. your tests are gonna be a good indicator about versions. Whenever you have external tests, it means your test is in another package
from the source code being tested, so the test has to use the code as if it's a 3rd party user. You don't have control over
internal things. So there's a good chance that if your write your test well and if your test starts to fail, that means that other
people who are using that package that's being tested, might have some code that works similar to the way that you tested and if
your test fails, there's a good chance that you're breaking the package for a future user, so it's a very good sign that you might need to 
bump the major version and make it clear to the users that ok this is a major version bump, it has breaking change to the API,
it's no longer gonna work the way it did, they have to tweak some stuff. So it'll also give you a good clue to whether or not
you need to do major version bump.
3. since the external tests have to use source code packages as if we were regular devs using those source packages, so we knew nothing
about the inner workings(because we don't have access to them) the tests tend to end up making much better examples and documentations for
other devs to check out. If we were using the private stuff of a package, another dev could be like: well I can't do that, so
I can't use your test as a useful example. But when we have to write our test from the external viewpoint(the external test),
we're using that package in the exact same way that another dev would, so they can take our test and use it as an example code.
We should be writing our examples as if it was an external thing. So it would be a useful example file.
4. one other reason(not specific to testing) is whenever you have to write your tests as an external package(external tests) and
you have to use your code as if you were a regular user(user of the package), it tends to lead to better designed code.
It means if you had access to all private stuff of the package, you might take shortcuts or cheats to write a quicker test.
But if you can't do it in the tests(external tests), you have to maybe spend some time building stuff and then you would find out
it's tedious, so you would realize you designed your package poorly, it's not user-friendly. So by writing external tests,
you get to beta-test it yourself. This is important because as devs, we often design a func so that it's easier to implement the func
and we will often forget about the end user whose calling that func, but in reality people probably gonna call that func
a lot more than you're gonna spend time writing it the first time, so we should generally design for the opposite. We should
try to make it easiest to **use** the func and maybe a little bit harder to actually make the func work with that input. Maybe we have
to do a bit more work to make the input into the right format, but this makes it easier for other devs to use it.
5. this point is not a benefit, it's sth that you **have to** do. The source code we wanna test, is in package `x`. Now this package
is being used by package y. Now we wanna test package x, but we also need package y. But we can't do that, because right now our test
is in package x and we would get a cyclic import. Instead, we can put the test into another package named `x_test`. Now there's no
cyclic dep because x_test is using both x and y but those two are not using x_test. There's no cycle anymore. That's another point
why you wanna use external tests. To fix the cyclic dep issue when you wanna write a test or example that requires you to import sth
that would eventually create a cyclic dep, but by putting it into an external package you can avoid that.

## 051 Exporting unexported vars funcs and types

## 23-023 Running tests in parallel
Why we might want to use parallel tests?

Sometimes running many tests in parallel can provide a ton of value.

Reasons why we may not want to use parallel tests(When to avoid parallel tests):

example use cases:
- simulating a **real-world scenario** - a web app with many users
- verify that a type is truly **thread-safe** - verify that your in-memory cache can handle multiple concurrent web requests using it

Note: Why a cache should be thread safe? Because it's **shared** across multiple threads. Why it's shared? Because we want to
calculate it once and use it multiple times, to avoid re-calculating it, because it's expensive. If it wasn't shared, it wasn't useful in the
first place! Because otherwise we had to construct and deconstruct it in every single req, so it wouldn't be useful.

So when building sth like cache or some type that needs to be thread-safe, it makes sense to setup the tests to run in parallel, to
verify that is truly thread-safe.

Parallelism is not free. It could mean more work and weird bugs:
- Tests can't use as many hard-corded values; eg unique email constraints. all tests are using the same email, now all of them is running,
  we will run into errors and tests would fail, because multiple reqs are using the same email, so we need to
  generate unique valid email addresses
- Tests might try to use shared resources incorrectly; eg image manipulation on the same image or sharing a DB that doesn't support multiple
  concurrent connections(like voltDB) or maybe the code being tested is writing to an output file.

When you run `go test -v`, in console we see:

`
=== RUN   TestA
=== PAUSE TestA
=== CONT  TestA
--- PASS: TestA
`

**CONT means continue which is used for paused parallel tests**.

This `PAUSE` occurs when we use t.Parallel(). **The test that have `t.Parallel()` get paused and they're added to a list
of tests that are supposed to run in parallel. Then go runs all non-parallel tests and after all of them ran,
it'll finally run the parallel tests.**

The reason for this is because go's testing package guarantees that only tests that have `t.Parallel()` will 
run in parallel together.

So by having `t.Parallel()`, they signal that only those tests should run in parallel.

So tests that don't have `t.Parallel()` won't run in parallel, so we don't have to do anything special for them to make sure
they won't run in parallel. All we have to do is **not** call `t.Parallel()`. Why we can't use a goroutine to run tests in parallel?
Because let's say we run TestA with a goroutine, so that it can run in parallel with others. But that won't stop other tests
that shouldn't run in parallel to run with TestA at the same time. Because there's no guarantee here. Other funcs will
be called at the same time with the one running in goroutine.

But `t.Parallel()` gives us that guarantee that only the tests marked as parallel, will run together and other tests that are not
marked, won't be running with parallel tests at the same time.

So when we have:

```go
package main

import "testing"

func TestA(t *testing.T) {
    t.Parallel()
}

func TestB(t *testing.T) {
  t.Parallel()
}

func TestC(t *testing.T) {
  t.Parallel()
}
```
When you run `go test -v`, output would be:

    === RUN   TestA
    === PAUSE TestA
    
    === RUN  TestB
    === PAUSE TestB

    === RUN  TestC
    --- PASS TestC

    === CONT  TestA
    --- PASS TestA

    === CONT  TestB
    --- PASS TestB

This means when go encounters a parallel test, it pauses them and puts them in a queue, run the non-parallel tests and then run
them by their order in the queue.

You shouldn't use parallel tests just for just speed boost! This is not a good reason. Because bugs caused by running tests
in parallel, you're gonna waste a lot of time trying to find out what's going on and debug it. That you're gonna loose all that
speed boost.
So use parallel tests for the use cases mentioned previously.

Most of the time your tests are gonna run fast anyway without running in parallel.

## 24-024 Parallel subtests
`24_test.go`

Subtests can also run in parallel, but they will only be run with other subtests from the same parent test in parallel.

The result is:
![](img/section-2/24-1.png)

Note that there is no t.Parallel in root of TestB.

So you see TestSomething and TestA get paused, then sub1 and sub2 subtests get paused, then sub1 and sub2 are continued(CONT) and
they both run and pass(or fails) before we ever get to TestA and TestSomething.
Note: TestB(parent of subtests) won't finish until it's subtests are finished running in parallel.

This is because sub1 and sub2 will run in parallel together(why together? because they share the same parent) and when they are finished,
all other parallel tests will continue to run.

So subtests run in parallel and after they finish, other test would run.

This behavior ends up with a benefit because it means we have a lot of control over which tests can and can not run in parallel together.

Let's say we want all of TestB and it's subtests to run in parallel with A and Something. For this, add t.Parallel() in TestB.
So all of the tests in that file will run in parallel.

![](img/section-2/24-2.png)

So all of the tests(subtests and others) will run in parallel.

If we only put t.Parallel() in  root of TestB and subtests not having t.Parallel():
![](img/section-2/24-3.png)

So it's important where to put t.Parallel().

## 25-025 Setup and teardown with parallel subtests
When testB runs, we get:
![](img/section-2/25-1.png)

Which is unexpected! Because teardown and even deferred teardown get run. And then parallel subtests continue and they finally are run.
But the teardowns are already occurred!

So if the setup included sth like opening up db connection or spinning up the app and defer were to close that db conn or app,
whenever the subtests get to run, they're not gonna have the db or app to use, they're already closed!

We wanna still queue those two subtests in parallel which is why `t.Run()`s has to return, otherwise go can't come down and start
next t.Run(). We need a way to tell us when all of the subtests have finished, so we can do the teardowns. Without that,
the t.Run()s will exit because they have t.Parallel() and we can't have a way of knowing when those tests are done, so we can start teardowns.

Solution: Wrap all parallel subtests in another t.Run() but this parent t.Run() doesn't have t.Parallel() itself.
Because the group shouldn't return(finish) until the entire group is finished running, because the group is not supposed to run in parallel,
it's supposed to run sequentially.

So whenever having setup and teardown for parallel subtests, you might have to wrap parallel subtests in another t.Run() to make sure
that teardown happening in the way expected.

## 26-026 Gotchas with closures and parallel tests
In earlier versions of go, when you run it, you see the test names are correct(since it's not in the callback), but the callback func of subtest is using 10 as i.
That's because we're not copying the value of i into the closure. The closure is just using the var i but that variable is incremented
everytime the iteration of loop executes.

The closures won't copy the value of i before they run and know to use that copied version. Now since all the subtests go into a queue,
i is gonna incremented to the last value before any of them subtests get to run.

**Solution: To fix any closure issue like this(typically in a for loop), copy the data into a local var and use that local var.
As long as you're declaring and initializing the local var before t.Parallel(), it's fine. But if this solution is done after
t.Parallel(), it won't work.**

Solution 2: Create another closure. This is not preferred.
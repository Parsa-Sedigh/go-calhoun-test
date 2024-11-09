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
## 25-025 Setup and teardown with parallel subtests
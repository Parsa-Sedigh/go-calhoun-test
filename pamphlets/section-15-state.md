# Section 15 - state

## 058 What is global state
A stateful program is one that remembers preceding events or interactions.

The info being remembered from these interactions is called the **state**.

Global state is info being remembered that is not isolated to your currently running process. Eg:
- package-level variables
- db entries
- it could even be in the form of struct fields in some code designs.

Some state can't be avoided - eg we may need a db in our app.

Generally speaking, global state can lead to unpredictable results and thus can make testing harder.

Global state is not recommended and we should get rid of it.

Global state: variables or some info that's being remembered that's not isolated to a specific goroutine. Like package-level vars,
struct individual fields or external stuff like db values which are global state that persist between every http req. Because multiple goroutines can all
be reading the same rows. But for example transactions can be used to get around this issue in dbs.

Global state should be avoided as we can, because it's hard to test with, it can lead to unpredictable results.

Example with pporf(github.com/joncalhoun/twg/pprof):

```go
package main

import "net/http"

func main() {
	http.ListenAndServe(":3000", nil)
}
```

Now visit `localhost:3000/debug/pprof`. Where did this all come from?

- [golang.org/src/net/http/pprof/pprof.go#L72]()
- [golang.org/src/net/http/server.go?s=73173:73217#L2391]()

uses a package global `servemux`.

Because of this, imported packages could end up altering how our program runs and similarly other tests could alter state and make tests fail.

- Tests can also lead to invalid states - eg if two parallel tests try to set states at the same time, they might lead to an incompatible state for
any tests. 
- Order also becomes important - if one test expects state from another test, it can lead to a specific testing order being required which
is often a bad sign. 
- State can also lead to very flaky or unpredictable tests which is bad - you want tests to pass if the code is good, and fail if there's an issue.

**Avoid global state if you can. We'll see how in the next section when we talk about dependency injection.** The rest of this section assumes
you can't avoid the global state and need techniques for testing with state.

## 059 Testing with global state (if you must)
Don't use global state, use DI instead.

Again, avoid global state if you can. Sometimes a test simply needs to interact with a db, or some other system that maintains state.
These tips are for those situations.

Let's assume we can't avoid global state.

If you must use it, here are a few ideas for making tests more reliable:
1. don't use parallel tests. They will very likely lead to issues. Because when we have global state, we don't want two or more things
all doing the same thing at the same time.
2. specify your test order more concretely if you need to(this goes hand-to-hand with parallel test). Note that other tests 
that don't use global state, can run in parallel. But still it's tricky to do parallel test, because we never know if some other
test that somebody writes and add t.Parallel(), might accidentally introduce sth that alters the state that the tests using global state,
are using. So if you're using global state, don't use parallel test altogether. But sometimes you can get away with it.

Note: Subtests can be run with t.Run() or using another func by passing t there like testThingA(t).

```go
package main

import "testing"

func TestApp(t *testing.T) {
	// inside here we can run specific tests in a specific order
	// subtests
	testThingA(t) // or t.Run(...)
	testThingB(t)
	testThingC(t)
}
```

You should be able to run subtests this way if you need them, but be careful about parallelism because adding any parallel subtests
could alter that order.

3. use separate setup/teardown between each test that alters state(we can do these between each set of tests).

```go
package main

import "testing"

func TestWidget(t *testing.T) {
	// setup - create stuff we need in the DB
	db := // open DB
	user := createTestUser(db)
	widget := createTestWidget(db)
	
	// run tests that use the data in the DB
	// ...
	testThingA(t) // note: maybe we have to make sure if testThingA alters the widget, maybe it needs to reset it back to what it was
	testThingB(t) // if testThingB deletes the widget, maybe it needs to re-create it once it's over.
	testThingC(t)
	
	// teardown - delete stuff that we made!
	resetDB(db)
}
```

This doesn't mean you necessarily need to create/delete entries between every test. You could setup the DB, run several subtests,
then tear it all down. Basically you need to decide how long the data should live for your specific tests.

Point (3) is, IMO, the most viable option but it requires you to NOT use parallel tests in many cases and you need to be aware of
what state you alter and reset it all.

Sometimes you can forget what state is altered or not even realize some state is altered which is another reason why testing with
global state is hard. Package-level vars are even harder, because we can't query them like a DB, for example serveMux var could get
a new route registered and we never see it. Maybe we expect some path to be 404, but it's not.


Note: Another approach is you could try creating unique things every time. For example, for tests that need db conn, you could write custom
code that creates whole new DB and every test uses it's own separate DB or create unique email address or user for every test, so they conflict with
each other. Remember that you generally shouldn't use parallel tests here, because it's still possible some other test in parallel could create the same user.

So that at the end of the test suite, we can wipe the whole db, we don't care which one anymore.
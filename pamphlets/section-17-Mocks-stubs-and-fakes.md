# Section 17 Mocks, stubs, and fakes

## 067 What is mocking
**Mocking: Replacing a real impl with a fake one that is intended to use in either development or testing. In other words,
mocking is passing fake implementation.**

Note: While mocking is really good in tests, there's still a benefit to use real integration tests where we actually use
concrete implementations instead of mocked ones and with this, we can verify things actually would work in prod. This is
beneficial because it's possible to write a unit test which uses a fake(mocked) obj like email client, passes, but then when we
go to prod, the mailgun client starts failing. The reasons could be:
1. the mocked obj and the real impl might be different
2. or we misunderstand how the real impl obj was intended to be used and our mocked obj didn't reflect that

So mocking can't replace integration tests. We shouldn't mock everything in our tests, we still should have some integration tests
where they use real impls.

## 068 Types of mock objects
Note: **Throughout the course, we use the term `mock` to loosely define any of the 5 terms below. We're gonna use the term `mocking`
for passing any of these 5 things to a func.**

Terms:
- Dummy: simplest of all these. It doesn't really exist in Go. Because a dummy is sth that you pass in which is not used at all.
For example, we're saying: Alright, this func requires an io.Reader, but we know it's actually not gonna be used in this specific test,
so we can just pass in a dummy. Now the reason we say dummy doesn't exist in go, is because if we know sth is not gonna be used,
we can just pass in `nil`. So there's no real reason to go about creating some sort of dummy obj.
- Stub: A stub is sth that returns bare-minimum. It's return value is dependent on what arguments are passed to it. Or maybe it doesn't care
what the inputs are, it always return a specific thing. We say: Alright, if they call this func with args `1, 2` , then 
that func should return this specific thing or do this specific thing. That func is stubbed.
So stub is not close to a real impl, **it's just stubbing out the behavior, so we can simulate sth in a test. When stubs can be useful?
Whenever we need to simulate very specific behavior.** An example for specific behavior: writing code for creating a file and we wanna
test for `ran out of disk` behavior. Well, testing this scenario is very hard because we don't want to fill up the hard drive just to
verify that case is handled correctly. So we can stub it with a func that always return that err.
- Fake(aka double): has a bit more sophisticated impl than dummy and stub. Usually it's not gonna be quite as intricate as the real implementation,
but it's gonna be sth that is comparable to the real impl. An example is you might have a SQL DB and you write some stores(repositories)
for it and they all use that sql db. Now for demo or testing or ..., you might occasionally want to replace that sql impl with an in-memory db.
So you could use a **fake** to replace your userStore(userRepository) and other aspects of the code.
**A fake has some sort of fake but almost realistic implementation.** **Stub vs fake: For example in UserStore.Create(), the stub
might return some successful response but it's not gonna necessary remember that user, so later if you wanna 
find that user(by calling UserStore.Find()), a stub might return a hard-coded user every single time, whereas a `fake` might keep track of
all the users you've created and might return those users because it's using some sort of in-memory store to keep track of them all.
So a fake is more realistic than stub.** A fake is much more useful for doing demos when you don't have all the code implemented, or
when you're doing tests where you need to create some users and simulate some of it, but for some reason, setting up a real DB might be impossible
or hard.

- Spy
- Mock

In go, spys and mocks are the same thing. These two could be stub or fake. The difference between stub and fake vs spy and mock is
stub and fake isn't gonna keep track of what methods are called, what args were passed or other usage info. A spy and mock keep track of
that info and they use that info. So they could say: Ok, I first had the Read() called with arg1 and then with arg2. You can later
use this info in your test to decide if the test failed or passed. 

We can use spys and mocks to track if a func was called with what arg and ... .

Q: Why in go spys and mocks are kinda the same thing?

A: In many Langs, you have things that write sth like: I'm expecting this implementation to get these 3 methods called and if they aren't
called, the test fails. **In those langs, a spy wouldn't fail if sth wasn't called the way we expected, but it'll give you info to look at but a mock 
might fail if you didn't call the func the way you expected.** In Go, there's not a good way to do this. So what we typically 
have in go is you'll have some mocked obj that you'll create or generate, and when you pass it in, whenever the func that
uses that mocked obj is done, you can check to see what was called but you manually have to say: I want this test to fail if it wasn't called like that.

Note: We're gonna look at mocks that are fakes and mocks that are stubs.

Note: It's important to remember that while spies and mocks do have some value in different tests(we can verify certain methods are
called or certain things happen), they don't replace integration tests entirely. So you could have a test passing with a mock or stub or
fake, but it might not pass with the real impl and that's where we have to make sure we're still using integration tests, but mocking in
general will allow us to write fast tests that we can mostly make sure are correct and then we can come back and run integration tests later
to verify mocks and stubs and ... , are having correct assumptions.

- others?

With all of these terms, there are two major considerations:
1. what type of implementation is provided?
2. does it track usage info? Like what methods are called, what arguments are passed into those method and ... .

```go
package main

import "io"

func Foo(r io.Reader) {
	r.Read(...)
}
```

Example of stub:

```go
package main

import (
	"os/exec"
	"testing"
)

func TestGitVersion(t *testing.T) {
	// here, we're stubbing the func that is assigned to execCommand variable. It's not a real impl.
	// We're simulating a very specific git version
	execCommand = func(name string, arg ...string) {
		return exec.Command("echo", "git version 2.17.1")
	}
	defer func() {
		execCommand = exec.Command
    }()
	
	want := "2.17.1"
	got := GitVersion()
	
	if got != want {
		t.Errorf("GitVersion() = %q; want %q", got, want)
    }
}
```

## 069 Why do we mock
Mocking is used for many different reasons.

We gotta evaluate: Is it worth it to not use the real impl in this test and to not have a 100% certainty that it's gonna work with the real thing?
Because anytime we use a mock, there's no guarantee that the mock is able to simulate what the real impl is gonna do. There's always gonna be
assumptions, because we're not using the real impl, we're using it's mock. So it's a tradeoff between not using the real impl in order to
being able to test sth. If we don't mock it, testing would be hard or impossible, but with this, we would loose the certainty that the 
real impl gonna work.

### Common situations we might wanna use mocking
Mocking is almost always done to simplify testing, or to make testing possible.

- setup/teardown may be simpler
- simulating a specific situation(often an error) can be easier
- external APIs may be slow or unreliable
- API may not have a test env, or it may have limits
- we can verify that specific behavior occurs - eg that we call the `EmailClient.Welcome()` after a user is successfully created

Let's look at some examples to understand this better.

### simulate specific situation
`twg/race_pass/users_test.go` - we use a fake implementation that wraps a real UserStore in `racyUserStore` in order to simulate a very specific race condition.
It actually has the real impl as the embedded struct, but we're using it to send a specific behavior, so this is one of the cases where our
mock is actually more complex than the real impl, because it includes the real impl and some other stuff.
Look at 29/race_pass folder.

Simulating a specific behavior or other types of errors that are hard to encounter in real impl, but you wanna make sure that your code handles it
correctly. You can use mocks for simulating these errs.

### setup/teardown
We saw how we can create a real DB and we can even seed it with real info. This is a good bit of setup/teardown in `twg/psql/users_test.go`.

Sometimes a test doesn't really need this all to actually test sth. Eg our signup code from earlier:

```go
package main

import "strings"

func Signup(name, email string, ec EmailClient, us *UserStore) error {
	email = strings.ToLower(email)
	user, err := us.Create(name, email)
	if err != nil {
		return err
    }
	
	if err := ec.Welcome(name, email); err != nil {
		return err
    }
	
	return nil
}
```
Setting up a real DB for this example might be overkill. Instead, it's easier to create a stub for this UserStore and to have it return
a fake user or err(whenever we want to) and then we can test both scenarios easily, we don't have to create a user that already exist
with duplicate email in order to simulate what happens when user wants to register with an already existing email. Mocks help us with this.

We can mock out the UserStore entirely and avoid any SQL setup - we just return a user when we want to test a successful situation and
return an err when we want to test an error case.

So mocks make testing specific behavior and errs easier and also the setup and seeding it and teardown it all down when 
you're done and ... is no longer necessary because we mock parts of it.

### External APIs
Email clients like we saw earlier are similar, but imagine you're using an API to order postage labels for your packages.

What happens if the shipping company doesn't offer a test API, so you can't actually run tests with that API integration? **We use mocking!**

Or a test env may exist but doesn't perfectly match production(happens more often that you'd guess).

Or test env has limitations and can't be hit as often as devs are hitting tests.
Or using a real API is just too slow and you don't want that many network calls. Or we want our tests to runnable without using internet
for using 3rd party APIs which require internet conn and we don't want the go test to sit there and halt because it doesn't have internet conn.

### We can verify that specific behavior occurs
Signup example again:
```go
package main

import "strings"

func Signup(name, email string, ec EmailClient, us *UserStore) error {
	email = strings.ToLower(email)
	user, err := us.Create(name, email)
	if err != nil {
		return err
	}

	if err := ec.Welcome(name, email); err != nil {
		return err
	}

	return nil
}
```

We might use mocking and an actual mock which is one that tracks which methods are called to verify that we actually call the
Welcome() method in some situations but DO NOT in others.

*Warning: The type of testing with mocks where we actually verify that specific method is called, is very close to
testing HOW a func implemented rather than that it does what it is supposed to do, so use the technique with caution.*

Most of the time when we write good tests, it's about verifying that the end result is what we wanted and we don't care about
how it happened. Example: If have a sort func, we don't necessarily care about the algo it's using(how it's implemented),
it's not important detail, instead in test we care about the result being sorted.

The problem with testing implementation details, anytime we make a change to the func, all of the tests fail and that starts
to make tests useless and the result is we stop caring about those tests(they don't provide much value anymore).

### Summary
We use mocks when the real thing just doesn't make practical sense. Whether it's gonna make the tests simpler or possible.

*Caution: Mocks can differ from the real thing so it can lead to tests passing with mocks but there is a bug. Integration tests
are still recommended, even with mocking.*

070 Third party packages
071 Faking APIs
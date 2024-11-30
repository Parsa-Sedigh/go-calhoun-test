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
Mocking is almost always done to simplify testing, or to make testing possible.

- setup/teardown may be simpler
- simulating a specific situation(often an error) can be easier
- external APIs may be slow or unreliable
- API may not have a test env, or it may have limits
- we can verify that specific behavior occurs - eg that we call the `EmailClient.Welcome()` after a user is successfully created

Let's look at some examples to understand this better.

### simulate specific situation

### setup/teardown


070 Third party packages
071 Faking APIs
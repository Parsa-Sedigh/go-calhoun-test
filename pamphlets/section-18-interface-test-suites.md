# Section 18: Interface test suites

## 072 What are interface test suites
They're a way to test multiple interface implementations all using the same set of tests.

Where this is useful?

When a func accepts an interface, a lot of the times we're expecting different implementations to be provided. For example a func
takes an io.Reader and it can be a file, a network conn.

Interface test suite allows us to write one single set of tests or a test suite and then anytime somebody writes an implementation for that
interface, they can pass it into that test suite and verify that it meets all the different needs of that interface.

In other words, interface test suites are just a test func(in a source file not a test file ending with _test.go), that accepts *testing.T
and whatever interface we wanna test(like `suite.UserStore`) and it allows us to run some set of tests that will verify that
the passed interface is implemented correctly(fulfills our expectations in the tests) and we can reuse this tes func for
any implementations of that interface.

This technique is useful for when:
- other people might be writing plugins for our code 
- or when we're using interfaces and we have multiple implementations - like a real(concrete) UserStore that uses SQL and a fake UserStore that
we use for testing. We can pass both of those impls to the test suite func. It verifies the fake version if close to the real version and it
shouldn't affect the tests because of being fake.

Useful for complex interfaces where the definition of the methods is not quite enough to verify that it's doing everything it's should.

For instance, you're building a service for managing storing files on cloud like aws s3, or blob store or google's file storage. People
can write different uploaders that handle uploading files to different services and then you could within your app, as long as uploader
implements the interface, it doesn't matter where we upload the files.

- `stub` dir is an implementation of the interface.
- `suitetest` dir is the interface test suite

Note: someone could impl the UserStore without being useful. So we might add a test suite for testing common use cases and make sure
that whenever you're implementing this UserStore, it's correct. Some common scenarios: after creating a user, we expect to lookup the user by
their id or email and it should return the same user. Or when we delete a user, subsequent calls to lookup the user by id or email
should fail.

For this, we create a suitetest dir or you can put it in the same directory of the source code, HashiCorp does this.
They would create a testing.go source file and then in the same dir(and therefore package) as the UserStore, they would add the testsuite.
But we won't do it like HashiCorp, we put the test suite in `suitetest` dir and therefore `suitetest` package. Either way is fine.

Note: The test suites won't be actual tests that's gonna run **automatically** by the test tools. Because the file name is the same
as the interface they wanna test. The filename doesn't have _test.go .

So since there is a `test` word in the package name, it's obvious that this is a testing utility.

---

Q: How the end user would use the test suite in their test? Because the testsuite packages are not gonna run by the testing tools automatically ... .

A: Look at stub dir. We have a stub impl of UserStore there. First we can make sure it implements UserStore interface. For this,
we create a test file named `userstore_test.go`. We can verify the stub impl actually implements the interface, by saying:
`var _ suite.UserStore = &stub.UserStore{}` in the test file. We can even put it in the source file!(userstore.go).
This line would fail if stub.UserStore doesn't impl the suite.UserStore interface. With this, we can check if our impl actually implements
an interface or not, **at the compile time**.

Whether you prefer to put this line in actual source file where you have the implementation, or to put it in a _test.go file,
is up to you.

Now in order to verify our impl passes the testsuite, we say:

```go
package stub_test

import "testing"

func TestUserStore(t *testing.T) {
	// stub.UserStore{} is our impl of suitetest.UserStore
	us := &stub.UserStore{}
	
	/* this func call will run the entire test suite. So we can completely test our impl of suitetest.UserStore, without writing much code. 
	Also if we wanted to write a sql impl or mongodb impl version, all of those are gonna have a test that looks a lot like this and would be short.*/
	suitetest.UserStore(t, us)
}
```

## 073 Interface test suite setup and teardown
We're gonna look at techniques to manage setup and teardown whenever each implementation of the interface might have it's own unique sets of
things that need setup and teardown.

### Approach 1
It's good to do the setup and teardown in test suite. So that when people are using that test suite, don't have to pass in to this
test suite, extra funcs or things to make that work. The test suite func can do this using defer funcs.
In the test suite func, we can put the setup and teardown next to each other.

### Approach 2
The first approach can be tricky sometimes. Because:

There's no guarantee that a UserStore(interface we're testing) might not have side effects that we don't know about. And maybe
the test suite func don't really care about this necessarily. As long as it impl the interface we're fine. But still it might be useful
to give the client that is using this test suite, a chance to handle teardown between test cases. For example in the case of testing
a UserStore in the test suite func, while Delete()ing that user in test might do the teardown, there could be some weird scenarios
where it's just easier to completely wipe the DB before we run a new series of tests.

So we need to provide a custom way to setup or teardown a test.

To support setup and teardown, there are multiple ways.

1. We could add more params to the test suite func(`UserStore` in `userstore.go`). Params like beforeEach, afterEach that are funcs.
Note: People could pass nil as those args, so we need to check them not being nil.
2. Create a `type UserStoreSuite struct{}`. It has fields for the actual interface we wanna test + beforeEach and afterEach funcs.
Then we can have a method called `All()` for running all the tests of the test suite. This way, if we you don't want to set extra
opts like beforeEach and ..., it's easier than having to pass nil. So if you start having to support a bunch of fields, you can create
a struct to pass those fields and have the test suite run via some method like `All()`.

## 074 Interface test suites in the wild

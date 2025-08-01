## 126-01. What to expect
Projects that we inherit are usually hard to test. So when we want to add tests, it'll be hard unless we refactor the code and then
when we wanna refactor, we can't refactor without having some tests to make sure we don't break things and we can't add tests because
the code isn't designed with being testable. So we're in a cycle where we can't do anything!

So when we inherit apps like this, it's gonna be hard to test. **We're gonna try to look at ways to incrementally add tests.**
While it's not gonna be perfect and has issues, it should give us much better approach than trying to refactor the entire project in one
big go. You can't refactor sth that's large and expect it to still work perfectly.

If you really need to test an app and you have no way of adding any tests, sometimes it's useful to either start with e2e tests or
manual testing. E2e test means instead of testing specific funcs or ..., we spin up a browser using chrome dp or gotty which is a
go package that spins up a browser and lets you write tests with that, or phantomjs, then spin up our app and do things that a user would and
verify we get what we expected. This approach is useful in short term until we get the refactor done enough until we can actually
move to some other testing technique. We can do this with e2e or manual testing. If e2e required a lot of effort compared to it's value,
use manual testing.

Because we could be using prod stripe key, so we would get charged.

---

If you consider to test the entire app, focus on leaf deps(deps & packages).

Note: If you have a large project with a lot of deps, to add tests and maybe refactoring those parts gradually:
Make a graph of all deps in your project and then start adding tests to leaf nodes. Because they're usually easier for starting vs
the main package or sth at the very root. And we gradually work our way up the tree until we get to the root node and by that point
everything else should already have tests and should be refactored, so it's more testable.

In swag, the db package is an easier place to start adding tests, because we know in go cyclic deps are not allowed. So `db` package
can't depend on the server. So db package should be testable by itself without server code being relevant.

The steps we're gonna take is similar:
1. add or improve tests the best you can
2. refactor the existing code into sth that's more testable or better written. So that we can add new features or ... .
3. repeat step 1. Since the code is more testable, can I improve the tests? Can I make them so they're less likely to break? And make
them more robust? Because a lot of times, the very first test that we add is when we're dealing with limitations of the code, but if we
keep repeating these steps, we'll incrementally make it better.

## 127 App overview
It's good to have an env var for ports because otherwise, by hardcoding a port for running the app or ... , that port might not be
available when testing.

To start refactoring & adding tests, we start at the leaf nodes like `db` package.

Why we start at leaf-level packages?

1. we can avoid writing big e2e tests because they're hard to set up and they're fragile and flaky
2. we could have hard-coded configs

## 128 Initial db tests
Note: Calling `panic()` inside `init()` won't terminate the program.

Note: The price col is of type int and it's in cents. So price set to 1000 means 10$.

Before making changes in db pkg, we wrote a test to make sure we don't break anything there as we refactor it. db_test.go is the test.

## 129 Creating the dbOpen function
Create the Open() func in db pkg.

Instead of format `host=x port=y ...` for psql conn, use `postgres://<user>:<pass>...`, because in the second one, people
can just provide **one** env var vs a bunch.

So **it's good to use a real DB in tests**.

## 130 What about mocks
When we're interacting with dbs in tests, it's more integration test not unit test.

go-sqlmock is a lib that is essentially a driver for sql where it provides a mock for all of that.

We need to test db as well because the mocks won't verify the query itself.
But remember that testing specific scenarios is hard when having real DBs. An example of where testing a scenario with real DB is hard,
is: What happens if we call GetCampaign one second before expiring or after? Well yeah we can stub the time part but still is hard.

We want integration tests anyway, if integration tests are enough, we're confident that the code won't break, so less need for unit tests.

## 131 Test harnesses and helpers
We can use fixtures or some sort of DB seeding for when we need some data to be in DB before every testcase.

---

One approach to clean the DB between testcases, is to run every testcase runs inside a tx, so it can rollback the entire tx after
the testcase is run, so we can continue other testcases. This is faster than wiping the whole db or redo the whole thing.
But the drawback is when we're running inside of a tx, sometimes things won't work exactly the same as they might in normal scenarios like
inside your own tx instead of frameworks tx that is run per testcase.

So using a tool for cleaning every testcase is probably overkill.

Let's set up a test harness(our own mini testing framework).

We're gonna set sth up that allows us to wipe the DB before every test and wipe it after every test just to make sure that everything is clear.

In a lot of tests we wanna do similar things like reset the DB. So we could write a func for this and call it in every testcase.
But we can also use test harness instead. But right now a test harness won't be a lot helpful, but later it will be.

For this, write a before or setup func(we named it setup).

Create `TestCampaigns` func.

## 132 Reviewing tests
Instead of having package-level funcs for db, create a type and move those funcs as methods.
The reason for this is, right now each func there uses a global db var without any dependency injection. By using methods
on a type that has all the deps, testing would be easier.

The only way we can still test the code that doesn't use dep injection is to overwrite global vars with the val we want in the test.
The drawback is we have to run all the tests one at a time, because otherwise a test might overwrite that global var and it would
become invalid for next tests(so we have to add code for cleaning up after each test).

Instead of all these hacks, create methods instead of package-level funcs and with this we utilize dep injection.

NOTE: If each individual testcase needs a set of dynamic data, we can use a func in the table driven tests to set up that data.

133 Testing specific times
134 First pass at refactoring the db pkg
135 Updating db tests
136 Testing the order flow
137 Extracting code for unit testing
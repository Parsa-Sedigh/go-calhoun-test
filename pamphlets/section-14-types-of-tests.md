# Section 14 Types of tests

## 053 Overview of test types
Terminology isn't necessary to write good tests, but can be useful in discussions and when thinking about the intent of a test.

In this section we'll briefly review the following types of tests:
- unit tests
- integration tests
- end-to-end tests

There won't be a lot of hands-on coding; instead we'll see made-up examples.

Projects and practice will be a much better place to really start to get a grasp on how these all can be used together to write a
robust test suite.

- Unit: testing very small things, like a func. Usually in isolation
- integration: testing at least 2 or more systems together.
- end-to-end: testing the entire application, or most of it. Usually in a way similar to how end users would use the app, but
that's not required. Sometimes there are e-to-e tests that might make API calls and verify those are correct, so it might not necessarily
load up the frontend app. But we still consider that an e-2-e test because it's an e-2-e test for the APIs, but it's not an e-2-e test
for the entire web app. So what we define our entire app can depend on how we view it. We might view it as 3 separate little apps that
all work together to make our web app, or it could be a single big app.

A common assumption with integration tests is that let's say we're testing the integration of system A with system b but we
assume that we can't change system B.

## 054 Unit tests
Testing very small things, like a func or a small type, in isolation.

For example, we can put some data into a type and test if it's methods work correctly.

Example:

```go
package main

import "testing"

// this is the unit - a function
func Magic(a, b int) int {
	return (a + b) * (a + b)
}

// this is the unit test
func TestMagic(t *testing.T) {
	got := Magic(1, 2)
	want := 9
	
	if got != want {
		t.Errorf("Magic() = %v, want %v", got, want)
    }
}
```

- Very common for these to require very little setup and to basically be `given x, do I get back y`?
- A lot of times, unit tests just test things of a single package.
- they're very easy to put in the format: `given x, do I get back y`
- they are run very fast, so we can run those everytime we hit save on the IDE to **get immediate feedback**

## 055 Integration tests
Test how at least 2 or more systems work together.

example:

```go
package main

import (
	"database/sql"
	"testing"
)

// in this code, we rely on an external system - the DB. If we wanna unit test it, we usually mock that external system

type UserStore struct {
	db *sql.DB
}

func (us *UserStore) Create(user *User) error {
	// ... this uses the us.db (the sql database) to create a new user entry from the user object passed in
}

/* integration tests might use a REAL database, meaning it's testing the integration of our UserStore with a real SQL DB and
not some mocked out DB.*/
func TestUserStore_Create(t *testing.T) {
	
}
```

In unit test, we would mock the external systems like DB. But in an integration test, we don't mock the external system and for example
we actually have a real DB.

When we have integration test, we have to set up things correctly, like connecting to db, make sure db has correct tables.

Why the separation of unit and integration tests matter?

Unit tests, especially ones with mocks, only test that another systems works as we expect it to work. Integration tests will verify that
our expectations of how the system should work are correct.

So we could have unit tests that mock dbs and report that the test passed, but since there was an error in sql query and it didn't
run in the unit test(since it was mocked), the integration test would fail because we use the real db and won't mock it in integration.

Put another way:
> unit tests do have one major disadvantage: even if the units work well in isolation, you don't know if they work well together.

https://testing.googleblog.com/2015/04/just-say-no-to-more-end-to-end-tests.html

Eg imagine you're using a payment API and you write the following code:
```go
package main

func main() {
	customer := api.GetCustomer(email)
	charge := api.CreateCharge(customer, 100)
	
	// now we assume the charge is successful and move on ...
	
	order := createOrder(...)
	shipment := createShipment(...)
	
	// and then we tell our warehouse to ship the item
}
```
In this case we might test our code with mocks(which we learn about later) by saying "if we call api.CreateCharge() then we can assume
the API will work and move on".

What happens if we didn't realize the API requires us to finalize a charge by calling another API before it actually deducts an amount from a user's balance?
CreateCharge() only creates an initial row in db, we need to call another API too.

Our unit tests, or any mocked tests are likely to pass, but integration tests would likely fail, because for example it would check the balance
from the db or ... . So with integration tests, we verify that our assumptions in the unit tests were correct.

Another example is when we're using stripe and we're passing wrong id(with the correct type) for an API. Now in the mocks, since we have mocked
the call to stripes api to always return successful res, we might think our source code is correct. But then if we go read the docs or
did integration tests which won't mock the stripe's api(in staging env of stripe OFC!), it would fail, because we're passing in
the wrong type of id. For example, instead of passing payment source id, we're passing customer id and that customer
doesn't have a default payment source yet.

**Usual mindset: I'm testing interactions between A and B, and can't change B.**

We can't change how postgres or stripe APIs or ... works, but we can change how our system interacts with it. So our code would be system A and postgres
or stripe is the system B. We can't change B, but we can change A and verify that it works with B correctly.

Sometimes you own both A and B. So B might be a package or service or sth that you own and run yourself. But for the purpose of
the integration test and from your perspective, you should assume that B is a working system and it's working as intended and we **can't**
change it. And this makes it easier when integration test start to fail, usually you know that the system B(system we can't change) is not
where the bugs are gonna be, **usually the bugs are gonna be in the system that you can change**. This is not always true though. But the
system that you're not supposed to change, should have it's own test suite, so if it does break or has a bug,
you should find that about it from another set of tests somewhere else.

- EX 1: I'm testing interactions between my PaymentService and Stripe, and I can't change the way Stripe's API behaves.
- EX 2: Testing your DB code with the `database/sql` package and a real DB.

Note: It depends on the lang, but in go, it's very common advice that when you're testing your UserStore(repositories) or sth that interacts
with the DB, instead of mocking the DB, it's almost always better just to use a real DB.

The reason this advice generally is good, is yes it will slow down your tests, but the fact that that integration test is verifying that it all
works correctly and it does what we expect, is enough of a benefit to us that it's worth that slow down.

So keep using the real DB when testing repositories, until you get to a point where it doesn't cale anymore and when you hit that point,
that's where we can look into maybe mocking out some of the slower tests and maybe have a separate test suite for slower integration tests
that you know are slower but you don't have to run all the time.

## 056 End-to-end tests
Tests the entire app, or most of it.

There can be a fuzzy line between integration and e2e. What happens if we test say 3 packages but not the whole app? 
Sounds like integration, but what if our app is only 4 packages total? TL;DR - the idea is more important than actually having
a concrete separation.

EX1: Code to ...
1. start your app
2. open up chrome
3. navigate to your app
4. login as a user would(entering data in forms)
5. and ...

- Great for simulating real user scenarios 
- Great for catching bugs - touches a ton of code

Could even involve MULTIPLE SYSTEMS; eg we could spin up an app with 3 DBs and while fake users interact with the system, we could
kill a DB to see that it all works as expected.

For example, in integration, we insert a user in DB and then start testing, but in e2e we start the app, open up chrome and use chrome to
create a user and login and ... .

### Cons
Not always great at pointing at WHY those bugs occurred or how to fix them(because it's perspective is from a real world user).
You can figure it out, but they don't tell you quickly and clearly.

Joke ex(both e2e & integration): Two doors both work in isolation, but put the two together and suddenly the doorknobs block each other
from opening!

## 057 Which test type should I use
Short version: A mixture of all of them. The exact ratio will vary from team to team and project to project.

Every type of test has a tradeoff. Figure out which tradeoffs make sense for your team/project. Do this with trial and error and/or
experience.

Why not just one type?

All unit means we don't test how a system works when put together. It would be like testing each piece of car - the tiers,
engine, ignition, etc - but not testing the final car once assembled.

Integration tests help make sure different parts work together - eg that steering works with the tires on the car - but 
doesn't verify the whole works.

e2e verifies the whole car works, but require us to have an entire car assembled and working to test it. It is also
less obvious what is broken if the car doesn't start.
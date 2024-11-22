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

## 056 End-to-end tests

## 057 Which test type should I use
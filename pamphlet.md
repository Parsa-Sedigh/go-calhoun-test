## 1-001 What is a test
## 2-002 Why do tests matter
### Why do tests matter?
1. Tests help find, fix and prevent mistakes(bugs, side effects, edge cases, etc)
2. Tests document expected behavior
3. Tests encourage us to write better code
4. Tests can speed up development

## 3-003 Writing great tests
### How do you write great tests?
What makes a great test is going to vary from project to project, situation to situation and team to team. But some general things that
always apply:

- test with a purpose(we should have some end goal)
- don't overdo it
- testing is a skill, skills evolve

## 4-004 Testing with a main package
## 5-005 Testing with Gos testing package
`sth_test.go` tells go testing tool that this is a test file. So that go tools will not use it for builds that are not tests but whenever
it's actually running your tests, it will use these files. So the code inside `sth_test.go` will never be included in your regular source code because
it's in _test.go source file.

This is not the format you wanna use for signaling failure:
```go
t.Errorf("Wanted 11 but received %d", sth)
```

## 6-006 What happens when we run go test
To see what's going on when running
```shell
go test
```

Create a test function and call `time.Sleep(time.Minute)` there and then run `go test`. Now while this is running, run this command which shows
any process that has `go` inside it's name:
```shell
ps -u <username> | grep go
```

You will see a process from a binary that is running from a temporary directory and the binary is named: `<package name>.test`.
After the tests are run, that temp directory will be removed.

What's happening behind the scenes is the go tool is building a new binary out of our source file that has our tests.

The main thing that the test tool is doing is it's gonna take all the <sth>_test.go source files, it's gonna use them to compile it to some sort of
binary(like our normal code) that's gonna run and then it's gonna store that binary in a temporary directory and then run it. Then at the end,
the go test tool will clean up that temp directory with the binaries in it.

## 7-007 File naming conventions
Anything with `_test.go` is gonna be a test source file and when we run `go build`, it won't include anything that is `_test.go`. This is a builtin behavior
in go tool.

- `export_test.go` to access unexported variables in external tests. Let's say you have an unexported variable in your source code file.
Now let's say you have a test that's not in the same package as the package that has that unexported variable. This is gonna make those variables
exported. So for testing from outside of that package, we can do it.
- `xxx_internal_test.go` for internal tests. 
- `example_xxx_test.go` for examples in isolated files. This is for isolated examples that can be used as test cases.

## 8-008 Function naming conventions
For testing functions, the name of the test function is: `TestXxx`

How do we test a method on a type?

A: Name of the test function should be: `TestType_Method`. So for example for testing the method Bark on type Dog, the test function's name would be:
`TestDog_Bark`.

You can also write multiple tests for a single function or method or type or ...  and for naming of test functions in this case, we just add
another _<xxx>. Notice the word after the new underscore is lowercase. For example: `TestDog_Bark_unmuzzled` or `TestDog_old`, `TestDog_young`,
`TestDog_puppy` and ... .

Since we have table-driven tests, we might not need multiple test functions for a function or method. So we can have just `TestDog` instead of
`TestDog_old` and cover all of the test cases for a dog by using subtests, table-driven tests and other techniques. Therefore `TestDog` is more
common than `TestDog_old` which tests very specific situations or use cases.

We're also gonna see how examples can be used as a test case and how they will show up in your documentation.

## 9-009 Variable naming conventions
How are we gonna name our variables in a testcase?

When calling the function you're testing, the result should be stored in a variable called `got` or `actual`.

Also we store whatever we **want** or we **expect**, in a variable called `want` or `expected`.

`got` and `want` tend to be used in go because they're short.

Also there is the name `arg` for the variable you pass to the function that is being tested and called.

```go
package main

func TestColor(t *testing.T) {
	arg := "blue"
	want := "#0000FF"
	got := Color(arg)
	 
	if got != want {
		t.Errorf("Color(%q) = %q; want %q", arg, got, want)
    }
}
}
```

## 10-010 Ways to signal test failure
Whenever we write test, we need a way to signal failure(to show that the test failed).

### Signaling failures
The `testing.T` type has a `Log` and `Logf` method. These work similar to `Print` and `Printf` in the `fmt` package. The only difference,
is that the `Log` and `Logf` are only gonna show up if the test fails. Now if you want to see those logs even if the test passes, you need to run
the `go test` command with `-v` flag.

Using fmt package for test functions is not good because it won't say in which test that log was printed out. So use Log and Logf whenever you want to
log sth in test functions for tracing info.

There are two basic ways to signal that a test has failed(mark the test as failed):
- Fail = fail, but keep running
- FailNow = fail now and stop test

You will rarely see these(`Fail` and `FailNow`) called though, because most of the time people will call the following methods which combine failures
and logging:
- Error = Log + Fail
- Errorf = Logf + Fail
- Fatal = Log + FailNow
- Fatalf = Logf + FailNow

Which do you use?
- If you can let a test keep running(it will keep giving you useful info even after a test fails), use Error/Errorf
- If a test is completely over and running further won't help at all, use Fatal/Fatalf. Because there's no point in wasting time running the rest
of the test.

If not using subtests, `Fatal` will prevent other tests in that function from running. Later when we learn about subtests, you will see
how `Fatal` becomes much easier to use exactly as you'd want/expect to work.

So if you aren't using subtests, for example you're running 6 different test cases without subtests(which is not recommended), Fatal and Fatalf will prevent
other tests from running.

## 11-011 When to use Error vs Fatal
```go
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()

	// we won't make a request to a webserver, we will just pass it to the handler. So just use an empty string as url.
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		/* If we couldn't event make the req to set up the entire test, we can' really run this test. So we want to stop the test entirely. */
		t.Fatalf("http.NewRequest() err = %s", err)
	}

	Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		// the rest of the test would be probably useless
		/* Note: We're not using the w and r which we passed to Handler() earlier, in the log here. Because they are complex objects and not suitable
		for logging. But you can pass sth to Handler() in the log string that makes sense. For example "GET, """.*/
		t.Fatalf("Handler() status = %d; want %d", resp.StatusCode, 200)
	}

	gotContentType := resp.Header.Get("Content-Type")
	wantContentType := "application/json"
	if gotContentType != wantContentType {
        t.Errorf("Handler() Content-Type = %q; want %q", gotContentType, wantContentType)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        /* If we fail to read the response body, looking at the rest of this test, we can't unmarshal it and ... all the code ahead is gonna be useless.
		So we use Fatalf here.
		
		Note: We didn't use `...; want = %s` because if we don't put any want, it's clear that the error was supposed to be nil. So writing:
		`; want = nil` is redundant, so we didn't write it.*/
		t.Fatalf("ioutil.ReadAll(resp.Body) err = %s; ", err)
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		// if we can't unmarshal, the rest of the test is gonna be useless
		t.Fatalf("json.Unmarshal(resp.Body) err = %s; ", err)
	}

	wantAge := 21
	if p.Age != wantAge {
		t.Errorf("person.Age = %d; want %d", p.Age, wantAge)
    }
	
	wantName := "Parsa"
	if p.Name != wantName {
		t.Errorf("person.Name = %q; want %q", p.Name, wantName)
    }
	
}
```

Note: We have to call `Header().Set()` before calling `Write` or `WriteHeader`.

## 12-012 Writing useful failure messages
General error message: `SomeFunc(%v) = %v; want %v`

Note: If the arguments you passed to the function that are being logged with Errorf or Fatalf, are small enough or are necessary, log them when writing
`SomeFunc()` at the above, so sth like: `SomeFunc(%v, %v) ...`.

Sometimes we add err to general err message, like: `SomeFunc(%v) err = %v` and a lot of times this stems from functions that return 2 values or
more values and the last value is an error and it wasn't nil.

Note: Whenever you have two return values from a function and the second one is an error, we assume the first one is the **useful** value and the second one
is the error that tells us if sth went wrong creating that useful value. So we don't put the name of the first value in the failure message, because
it's implied(that's the useful value). So instead of this: `t.Fatalf("SomeFunc() x = %s")`, we write: `t.Fatalf("SomeFunc() = %s")`.
But with the err it's less obvious. Because it's not the first argument, so we write it in the failure message: `SomeFunc(%v) err = %s`.
If you have a function that for example returns 3 things or it's less clear which is the primary(**useful**) thing it's supposed to be returning,
we need to be clear when you're writing the failure messages is to write what was what? So write the name of the argument as well.

Tip: This tip is sometimes different in other languages or frameworks: THe common order for `got` and `want`(`actual`, `expected`) is:
```go
if got != want { // actual != expected
	
}
```
The first part is always what we got. The **second** part is always whatever we **wanted**.

The same thing goes for general failure messages: `t.Fatalf("SomeFunc() x = %s")`. So whenever we print the error message, the order is:
1) function we called - `SomeFunc(%v)`
2) what we got - `err = %v` or `%v;`
3) what we wanted - ` want %v`

Testify encourages the opposite order.

EX) Maybe p has more fields but we're interested in two fields, so we only wrote those two fields in the string passed to `Fatalf`.
```go
err := SomeFunc(p)
if err !=nil {
	t.Fatalf("SomeFunc(name=%s, age=%d)", p.Name, p.Age)
}
```
So don't worry about the fact that you might not be printing exactly what was passed in to `SomeFUnc`. Print the ones that are useful. Because
maybe the person struct would be too big. Doing that will make debugging hard.

```go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestPick(t *testing.T) {
	seed := rand.NewSource(time.Now().UnixNano())
	t.Logf("Seed: %d", seed) // either log the seed here or in Errorf below
	r := rand.New(seed)
	arg := make([]int, 5)
	
	for i:=0; i < 5; i++ {
		arg[i] = r.Int()
    }
	
	// Pick is the function that we're testing obviously and it always returns the third element of an slice
	got := Pick(arg)
	if got != arg[2] {
		t.Errorf("Pick(%v) = %d; want %d", arg, got, arg[2]) // or: t.Errorf("Pick(seed=%v) = %d; not in slice", arg, got)
    }
}
```
The argument passed to Pick() in the Errorf message could be long. So the err message could be just: `t.Errorf("Pick(seed=%d) = %d; not in slice", seed, got)` 

Note: It's a good to log the seed. Because whenever it fails, we could see what the seed was, so we can sorta repeat or reproduce it. We could assign
the seed used in failed test to the `seed` variable as a hard coded value and run the test again(seed := 156726828).

**This was how we reproduce random number bugs in tests.**

We can use `Fatal` or `Fatalf` more frequently whenever using subtests because ending a subtest immediately doesn't stop the rest of the subtest from running.

## 13-013 A basic example as a test case
Go's testing package allows us to use examples as test cases. Remember that a big part of testing is helping people understand how the code should be
used(what this code is expected to be doing).

We can see the examples in the documentation, but we want them to serve as a test as well.

Another benefit of using examples as test cases is that the example code gets out of date, it gets stale and it don't work correctly. But by writing
these as test cases and they're actually being run every time we test our code, we have this guarantee that our example is always uptodate and working because
it's always being run as a test case.

Examples are exactly the same as tests with 2 minor details:
1) The function name starts with the word `Example` instead of `Test`
2) No args are passed in

The second part leads to how we actually test it? Because we don't have access to `*testing.T` .

A: We use a special comment that has `Output:` in it to tell the testing tool what we expect this example to output.

```go
package main

import "fmt"

func ExampleHello() {
	// Hello is the function that we're writing the example and test for
	greeting := Hello("Jon")
	fmt.Println(greeting)
	
	// Output: Hello, Jon
}
```
Now run `go test`.

Note: You can the expected output on a separate line, so:
```go
// Output:
// Hello, Jon
```

You also see that the `fmt.Println(greeting)` is actually not being printed when we run the `go test`. The example function has this special behavior.

## 14-014 Viewing examples in the docs
```shell
godoc -http=:8080
```

Now by visiting localhost:8080 :
1) you can always see the go docs even without internet.
2) go to packages page and there you can see your packages and their docs

godoc uses the name of the example functions to decide where to put that example. For example if you have `ExampleHello`, it will put that function under
the type `Hello` in the docs.

So naming conventions are important for making the examples go to the right place.

If you name a function `Example`, it will be a **package-level** example instead of being an example for a specific type.

So in the overview of package in go docs, under `Examples`, you'll see `Package` meaning it's a package-level example.

If you write an example named: `ExampleHello_spanish`(note spanish has lower case s), in the go docs, it will show up like: `Hello (Spanish)` and
if you click on it, it would be under `Hello` with the name: `Example (Spanish)`, meaning it's spanish example, even though that's not what it is, but
that's what it's showing you. So it's using the second part(after underscore) to create a specific name for the example.

## 15-015 Unordered example output
Q: How we deal with examples where the output is not in a specific order? For example if we're dealing with channels, there's no guarantee
for what order we might get the message depending on if we're using goroutines. Or we're ranging on a map. Sometimes we can't guarantee
the order, that's just sth we don't have control of.

Now we can do some logic to rearrange them, but this would make the example complex.

```go
package main

import "fmt"

func Page(checkIns map[string]bool) {
	for name, checkIn := range checkIns {
		if !checkIn {
			fmt.Printf("Paging %s; please see the front desk to check in", name)
		}
	}
}

func ExamplePage() {
	checkIns := map[string]bool{
		"Bob":   true,
		"Alice": false,
		"Eve":   false,
		"Parsa": true,
	}

	Page(checkIns)

	// Output:
	// Paging Bob; please see the front desk to check in.
	// Paging Alice; please see the front desk to check in.
}
```
Now run `go test`. This will fail because the order is not deterministic. You could even have the test PASS and not realize the issue.

So the question is: **How do we test this without caring about the order?**

We use `Unordered Ouput:` instead of `Output:`.

## 16-016 Complex examples
`Example_crop` is a package-level example.

We can get away with panic()s because this is an example and we know it's gonna stop the example from running and print the wrong output and example
should fail.

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func Example_crop() {
	/* Note: You can have a base64 representation of an image hardcoded here to be able to test these and then you can do:
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img))*/

	var r io.Reader

	img, err := Decode(r)
	if err != nil {
		panic(err)
	}

	err = Crop(img, 0, 0, 20, 20)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("out.jpg")
	if err != nil {
		panic(err)
	}

	err = Encode(img, "jpg", out)
	if err != nil {
		panic(err)
	}

	fmt.Println("See out.jpg for the cropped image.")

	// Output:
	// See out.jpg for the cropped image.
}
```

**Note:** The import statements are not shown in the examples by default. How we can show them? Because maybe we have some important import statement that
developers could forgot to import the way we intended, like doing empty imports because of **initialization side-effects**?

**Q:** How can we show side-effects that happens by doing empty imports like: `import _ "image/png"` (assuming there's an `init` function in that imported
pacakge that does some **initialization side-effects**) in our examples? 

The way you can do this in go is instead of putting that example function in a test file, put it in the source file. So move `Example_crop` to
`example_crop_test.go` and in a different package called `example_test`(just another package with an arbitrary name) . The naming convention
here is: `example_<test case>_test.go` or we could also name it: `<test case>_example_test.go`. In these files, there should be only a **single** 
example function. So we would only have `example_crop` in that file in this case.

Even the standard library is not consistent about naming conventions. The former is better where you have `example` as the first word.

The rules:
1) the test file has to have `_test.go` like `example_<test case>_test.go`
2) there can only be a single example function inside this file
3) it has at least one other function, type, variable or constant. If there isn't, just add a dummy thing like a dummy variable in that file, though
typically when you have a file like this, you would use that other thing in that file

By doing these, in the godoc you will see that the example shows the **entire** source file(including package name at the top and import statements).

```go
package example_test

import (
	
	// Needed for initialize side effect
	_ "image/png"
)

// typically you would have one and only one thing real here instead of this dummy variable
var file string = "this is not used"

func Example_crop() {
	// ...
}
```

## 17-017 Examples in the standard library

## 18-018 Table driven tests
It's common that we have multiple tests that use very similar setup or require the same checks or share some common functionality throughout the test.

We're gonna learn different ways for testing multiple test cases within the same function.

Look at `underscore.go`

To have multiple test cases, we can have multiple test functions, like: `TestCamel_simple`, `TestCamel_spaces`. But a lot of code between these
functions is gonna be similar Instead, we should use table-driven tests which is just a design pattern.

In the struct for testcases of table-driven tests, we define whatever arguments we need for every testcase.

When ranging over testcases struct, we can name every element `tc` or `tt`.

We should use Errorf when ranging over testCases in table-driven tests, because if you use Fatal, the rest of the testcases in a table won't 
run(in subtests we can use Fatal).

## 19-019 Generating table driven test code
In vscode, highlight a function and right click > `Go: Generate Unit Tests For Function`

It'll generate a test function and the test table would have a field named `name` which specifies the name of the test. If the function
we're testing only accepts an argument with the type string, you can use that argument as the name of the test as well.

We use `t.Run()` to initiate a subtest.

You can use `github.com/cweill/gotests` to generate table-driven tests in your IDE as a plugin(it can be used from command line as well but it's not
convenient).

This is another reason we use table-driven tests(even if there's only one test case) because we can generate them automatically using these tools.

## 20-020 Subtests
You can pass a named function to the t.Run() like:

```go
package main

import "testing"

func TestSth(t *testing.T) {
	t.Run("some name", someSubtest)
	t.Run("app test", appTest(app))
}

func someSubtest(t *testing.T) {}

func appTest(app *App) func(t *testing.T) {
	return func (t *testing.T) {
		// you can use the app here
    }
}
```
As you can see `someSubtest` doesn't have to follow normal rules of starting with the word `Test` because we're kicking it off **ourselves**. You only
need to have a test start with `Test` if you want the testing package to kick it off.

Note: You can have a test function with the name `testUserStore_Find` but go tools will give a warning about that that says:
`don't use underscores in Go names; func testUserStore_Find should be testUserStoreFind`. But it's good to keep the name of outsourced test functions that are
supposed to be closures when running t.Run() , like normal test functions, even though they're not exported.

AGAIN: Why we had these kind of functions? Because maybe we wanna pass some stuff to them. So we create another function with normal name of test functions
but unexported and they will accept the additional things we want them to have access, but they return a function that gets *testing.T . Yeah
we could do this with closures too. Example:

```go
package main

import "testing"

func testUserStore_Find(us *UserStore) func(t *testing.T)
```

If you have a shared setup and shared teardown across a bunch of tests when running subtests, you would pass a named function to `t.Run()` . 

You can nest `t.Run()`s.

Note: Instead of using a slice of structs(`[]struct{...}`) as test table, you can use a map of strings to structs(`map[string]struct{...}`) and in the
second case, the `name` field is gonna get removed from the struct and gonna be the key of the map.

See `TestCamelMap` func.

Another cool thing about using subtests with table-driven tests is we can fail our tests and this won't make other testcases in the table to not run.
**So if you use Fatal or Fatalf in a subtest that is being run with a table-driven test, that won't stop other testcases.**

So when having subtests, we have a fine-tuned control like when a test should end. Also subtests give us control over which
tests run in parallel of each other vs which parallel tests can have shared setup.

## 21-021 Shared setup and teardown
This is why we generally wanna use subtests.

To run a test you need to set sth up(you need to do sth to get that test ready to run).

Instead of using a fake interface for sql, we can set up a test database and use that for tests. Because it makes sense. You wanna make sure
this whole thing works with a real DB, so why not do that in tests? Therefore, we would have a shared setup and teardown.

For this you can create a test database and name it like: `test_...`.

`github.com/lib/pq`

If we can't set up the DB correctly, we can't test it reliably, so we can use Fatalf when setting it up.

Note: You can do the teardown in a `defer` statement or you can put it at the end, but using defer is preferred, because we could forgot to write it
at the end of the function.

We can create a helper function for setting up things(like db), but even with that, we want to run it once and then run a bunch of tests
using that setup.

## 22-022 TestMain
## 23-023 Running tests in parallel
## 24-024 Parallel subtests
## 25-025 Setup and teardown with parallel subtests
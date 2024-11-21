# Section 12: Additional testing flags

## 043 Benchmarks
The name of the benchmark funcs start with `Benchmark` so like `BenchmarkX`.

**Note: Don't pass b.N into other funcs.** It's not the number of ops. b.N is how many times your run your benchmark. The actual benchmark
code that is inside of `for i :=0; i<b.N;i++`, should be the same every time. So for example if you have a count func that counts to some number,
do not pass b.N to that func. Pass some static value like 100 to it. The reason for this is b.N is gonna be changed until the testing tool
think that it has a consistent reading for how long this op takes. But if the op(like count() func) is changing in how 
long it takes(because we're passing different variables(like b.N) as params of the func being benchmarked), it's gonna be hard for the
testing tool to get an accurate measurement. So do not use b.N as vars or params. So it would be:

```go
package main

import "testing"

func BenchmarkX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// do not use b.N here
		count(100)
    }
}
```

Why we need b.N?

Because there's all sorts of things that could alter how fast the code runs. For example(let's say we're benchmarking fmt.Sprintf()), 
if some other program is running and trying to print to the screen. Or if your entire computer has some I/O blocking op for some reason that
might slow down the benchmark or ... . But if we run the code enough times(b.N times) like 10000 times, we should start to get consistent
average of how long that code takes to run. So b.N is used to figure out what that average should be.

```shell
go test -bench .
```

Note: We can put the hardcoded values that we pass to the funcs being benchmarked, as the name of the benchmark func.
So:

```go
package main

import "testing"

func BenchmarkFibRecursive5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibRecursive(5)
    }
}
```

Note: we can write a generic benchmark func for our benchmarks of a func, look at `benchmarkFib` func.

Benchmarks can help finding out that new changes made the code slower, so we should fix it and ... .

## 044 Verbose testing
When running tests, using -v flag makes `t.Log()` to appear. If the test fails, `t.Log()`s will appear regardless of using -v.

```go
package main

import "testing"

func TestSth(t *testing.T) {
	// detect if we're in verbose mode
	if testing.Verbose() {
		// ...
    }
}
```

## 045 Code coverage
To use the coverage tool:
```shell
go test -cover
# result would be sth like:
# PASS (meaning our test passed)
# coverage: 25.0% of statements
```
How it calculates the coverage?

Behind the scenes, the go tool will:
1. take all the lines and add a counter, so it can tell whether or not each line of code ran
2. And then it will run tests which will execute those lines. 
3. And then at the end it will check which counters ran vs which counters never got incremented and it will calc a percentage based on this

So each line will get a counter of some sort. It won't count each line more than once.

In addition to using the coverage tool to see the **coverage percentage**, we can also get a coverage profile:
```shell
go test -coverprofile=cover.out
```
It generates cover.out file. And with that file, we can do:
```shell
# -func tells you for every function, what the coverage is
go tool cover -func=cover.out

# this generates an html page that shows which lines were covered, but making them green and the ones that aren't covered in red
go tool cover -html=cover.out
```

### What coverage is good for?
Code coverage is great for finding out whether we're writing enough tests.

### Bad usages of code coverage
Note: Having 100% code coverage doesn't necessarily indicate that you're catching 100% of your bugs. Having 75% code coverage
doesn't mean 25% of the bugs are gonna slip through.

hard requirement for having code coverage higher than certain number: 
Do not force people into having certain percentage of code coverage. It can lead to bad behavior.
Why? Because when you have a strict requirement of some code coverage, deleting code becomes sth that people might not want to do
because that would make the total code coverage less.

This requirement also forces people to write more lines of code by splitting things into multiple times like
breaking things up into more variables, in order to get more code coverage.

### What do we prefer instead of code coverage requirements?
We need to have a guideline. For example our guideline is 80% of our code needs to be covered with tests.

Now if a dev submits a PR and that new code is covered less than 80%, not that by merging it, it would cause the overall codebase to go
below 80%, our whole code base could have less than 80% code coverage, but new PRs should be above that line.

We'd like to tell him: at least explain why you're not complying with our guideline?

**So code coverage is useful as a guideline, not as a rule.** By being a guideline, we can skip it where it's need to be skipped.

## 046 The timeout flag
## 047 Parallel testing flags
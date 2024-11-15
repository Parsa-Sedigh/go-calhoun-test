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
## 046 The timeout flag
## 047 Parallel testing flags
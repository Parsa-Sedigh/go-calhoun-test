## 034 Building things with helper functions
### General helper funcs(not just for comparisons)
Helper funcs are not just for comparisons. We're gonna look at other use cases that make sense to write helper funcs, like
helper funcs for **building** things.

### Different ways to build things
**Since we open up db conns or ... in the builder funcs and we can't put the teardown logic directly in those builder funcs(because
they would be done fast and the teardown logic would be called before actual tests are finished), a common pattern with builder funcs
is to return a func from them and we put the cleanup and teardown stuff in that returned func.**

The helper funcs as builders, can accept `*testing.T`. They can use it for calling things like `t.Fatalf()` or `t.Helper()`.
By accepting `*testing.T` it implies that the func is a helper func and also we can fail the test in that builder func. By passing
`*testing.T`, we can fail the err in that builder func, so it doesn't need to return an error anymore. We can fail there.

Look at `userStore` helper builder func:
1. it accepts `*testing.T`
2. it returns a teardown closure. So while in that returning closure, we put a lot of logic, we don't expose them to the user of
the func.

**You can put the test helper func into their own package and name that package `xtest` like `iotest` in std lib,
or put them beside the source code that they're testing.**

Test helper funcs have 4 categories:
1. helper funcs for comparison
2. helper funcs for building stuff
3. helper funcs for generating data for tests
4. helper funcs for validating data

## 035 Generating test data
Helper funcs for generating data.

## 036 Go's testing quick package
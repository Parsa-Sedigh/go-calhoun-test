# Section 11 - Controlling which tests are run
## 038 Running specific tests
For example, we have integration tests that don't need to be run for some reason like we don't have a 3rd party dep installed or ...,
so we wanna skip them.

```shell
go test -v -run TestGotcha/arg=3

 # sub2 is a subtest in the `group` of TestB. Note that $ means `nothing after here`. So it won't run group2 or ..., only `group`.
go test -v -run TestB/group$/sub2
```

## 039 Running tests for subpackages
Maybe you wanna only run the tests for specified packages, like packages inside of a directory.

```shell
for pkg in *; do go test "./$pkg"; done
```

## 040 Skipping tests
For example you only wanna run your unit tests and not integration and end-to-end tests. Or short tests and not long tests.
Or some tests that need certain tools to be installed and you don't have it installed, you wanna skip it.

If `-short` flag is not set, testing.Short() returns false.
```shell
# with this flag, testing.Short() returns true.
go test -v -short
```

Note: When we run `go test`, intuitively it means we wanna run the short ones. So in unit test, we don't need to check for anything.
In other words, unit tests run by default.

In long tests, add a check that if `testing.Short()` is true, we wanna skip them. So that flag needs to be set explicitly,
so this approach has shortcomings.

To skip tests, we have 2 approaches(aside -short flag which is not a scalable solution):
- custom flags
- build tags

## 041 Custom flags
Note: TestMain runs before all of our tests run.

## 042 Build tags
With build tags, we don't need to add if blocks unlike custom flags, to check whether we can run the code or not.

if test a file has `+build=psql` at the beginning, to run it, we need to set the psql build tag:  
```shell
go test -v -tags=psql

# to pass multiple build tags, we need a space separated list but inside quotes, so that they all get passed in as one arg:
go test -v -tags="psql mysql"

go test -v -tags="integration e2e"
```

Note: By having +build integration, it means all the source code for that file are not included in the build if we don't add the build tag
when building the binary.

So build tags are better than custom flags because of:
1. not having to write if blocks
2. the files won't even get into the build if the related build tags are not included
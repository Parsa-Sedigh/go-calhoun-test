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
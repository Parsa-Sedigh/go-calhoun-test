# Section 17 Mocks, stubs, and fakes

## 067 What is mocking
Mocking: Replacing a real impl with a fake one that is intended to use in either development or testing.

Note: While mocking is really good in tests, there's still a benefit to use real integration tests where we actually use
concrete implementations instead of mocked ones and with this, we can verify things actually would work in prod. This is
beneficial because it's possible to write a unit test which uses a fake(mocked) obj like email client, passes, but then when we
go to prod, the mailgun client starts failing. The reasons could be:
1. the mocked obj and the real impl might be different
2. or we misunderstand how the real impl obj was intended to be used and our mocked obj didn't reflect that

So mocking can't replace integration tests. We shouldn't mock everything in our tests, we still should have some integration tests
where they use real impls.

## 068 Types of mock objects


069 Why do we mock
070 Third party packages
071 Faking APIs
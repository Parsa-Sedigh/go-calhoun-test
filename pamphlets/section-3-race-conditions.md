## 027 What is a race condition
### Testing race conditions
prerequisites:
- goroutines
- familiarity with waitgroups
- familiarity with interfaces

### What is a race condition?
![](img/section-3/27-1.png)

Note: We could add `b = balance` before the line that updates the shared resource(balance var) which is this line: `b -= amount`.
But that wouldn't fix it altogether because scheduler might decide to pause the current goroutine. So even though the balance got
refreshed, current goroutine couldn't use that refreshed balance and the other goroutine might change it and then when the
scheduler continues our goroutine, we get a changed balance and then update it which gonna update it to a wrong value again,
although in prev line we refreshed it.
`go
time.sleep(time.Second)
b = balance
// ...
b -= amount
balance = b
`

Solution: We need to use synchronization techniques like channels, mutexes.

## 028 The race detection flag
```bash
go test -race
go run -race thing.go
go build -race thing.go
go install -race pkg
go get -race golang.org/x/blog/support/racy
racy
```

Q: Why we might wanna use -race with build cmd?

A: Maybe we're building the binary for 5 different prod servers and we wanna check for race conditions. We can build **one** of the binaries
with -race on the `go build`, so one of the prod servers will get this binary build with -race flag. So maybe 20% of the traffic
in the prod(if they get equal traffic) will go to that binary. So we can check the logs and see if we get the data race warnings in prod or not.

This is the simplest option, but doesn't always catch race conditions.

Can even run a subset of your prod with this flag enabled to watch for unknown race conditions.

**Race flag(-race) doesn't always catch the race conditions.** For example, if the race involves a DB read/write scenario 
the "race" isn't in memory in go code, it is in how we interact with the DB.

See github.com/joncalhoun/twg/race_fail to see this in action.

The race_fail folder has a race condition and it causes the test to fail, but running go test -race won't report that race.

## 029 Testing explicitly for race conditions
By wrapping the UserStore and making the Spend() func to accept an interface, we can test for a specific race condition. Note that
Spend() func could only be accepting a concrete UserStore struct and not an interface, but that's not testable, including
testing for race conditions.

So having interfaces is crucial for testing, to pass in whatever we need to.
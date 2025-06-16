## 084 Why are dates and times problematic
Time & dates are hard to get just right because there are a million edge cases.
- timezones(hours, 30m, 45m, etc - not just 1hr offsets!)
- daylight savings time
- travel(switching zones for users)
- leap year(extra day in Feb)
- one-time adjustments to adjust for drift
- sometimes entire days are skipped.
  For example: Somoa skipped Dec 30 in 2011 - went from 29 to 31 in order to adjust the international date line location to
  make doing business easier
- feb 30 has existed

Note: One solution is using monotonic clocks and go has them.

Monotonic clocks always give us the right time. For example, if an hour has passed even if the daylight savings rolled the clock back,
the monotonic clock still gives us 1 hour.

Note: Timeouts are deterministic for most part. But we might have a computer that does the work in 30s, but another one that does it in
2m. If you set the timeout as 1m, one computer doesn't do that work.

So there's always gonna be this situation where we have some threshold where one computer will finish sth before timeout and another one won't
and this could affect your test. Especially if you have a test to simulate sth to happen before the timeout or you want the timeout to
actually expire.

For example, we wanna test sth that runs every hour. We don't want to actually wait for an hour for the test to verify it works.
Solution: We wanna simulate time.Sleep() and timeouts to have more control over them.
```go
package main

import (
  "context"
  "time"
)

func DoStuff() {
  for {
    // do stuff

    // then sleep for an hour and do it again!
    time.Sleep(1 * time.Hour)
  }
}
```

## 085 Inject your time and sleep functions
Look at `timing` folder.

The simplest way to make timing funcs easier to test is to inject those time-related funcs that you're gonna use, by using dependency injection.
With this, we can override the `Sleep` and ... funcs.

So with pretty much all these timing funcs, a good approach is to use dependency injection to be able to customize(like mocking) the funcs.

### Handling time.Now() in tests

### Handling time.Sleep() in tests


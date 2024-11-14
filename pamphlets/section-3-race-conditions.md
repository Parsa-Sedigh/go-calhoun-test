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
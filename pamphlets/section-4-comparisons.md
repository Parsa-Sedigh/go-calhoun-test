## 030 Simple comparisons
Note: While in many languages, if two objects that have the same fields but have different addresses, they're not equal if they're
compared using == operator. That's because in those langs they're trying to differentiate between checking if memory addresses are
equal(using ==) vs checking if underlying values are equal.

But in go that's not the case. Whenever you have two structs, if they have same fields and values, they are equal, even though
they're memory addresses are different.

**Note: Structs containing functions as their fields can't be compared.** So you can't compare structs like these:
```go
package main

type Dog struct {
	// ...
	bark func() // this struct can't be compared with other structs because of this field
}
```

Note: You can get the memory location of a var using fmt.Print(**%p**, &x).

Pointers are also comparable. If they point to the same variable(if they point to the same memory location), they're equal.
Also if both are nil, they're equal. Note that go doesn't look at the memory location of the pointer variable. It looks at where it **points**.

## 031 Reflects DeepEqual function
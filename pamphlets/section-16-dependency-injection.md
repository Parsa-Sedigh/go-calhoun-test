## 060 What is dependency injection DI
DI is a design pattern.

At the beginning, DI seems like a fancy way of saying passing arguments(dependencies we need), but this simple pattern
is gonna enable us write better code in large projects. DI enables more implementation agnostic code.

## 061 DI enables implementation agnostic code
Still there's a problem with DemoV2. It requires the caller to pass an EXACT struct(logger) - it's very strict in what it expects.

One way to make it better, is to look at what stuff we're using from logger and be more generous about what we expect.
One way to do this is to expect a func. Look at `DemoV3`. `DemoV3` no longer cares if we're using any specific package.
As long as it gets a func that it can use for logging, it doesn't care. So the param is implementation agnostic.

**implementation agnostic is achieved through interfaces.**

Another upside of using interfaces over concrete types is, when we expect concrete types, it's not clear what we wanna use
from that type and this makes it harder to use this func that accepts the concrete type. Because the func could do anything
with that type. This is especially true when the concrete type is complex and has a lot of fields and are hard to set up, we don't need which fields
are used in that func. Look at DemoV2, it expects a concrete type which is not good.

The upside of `DemoV5` param, is we don't have to pass Logger interface to every single method call.

## 062 DI makes testing easier
How being implementation agnostic leads to easier testing using DI?

It gives us the ability to replace the implementation that we use in our source code with an impl that still does the things
in realistic way, but we can customize it to test for specific things(they are more test-friendly).

063 DI and useful zero values
With DI, we need to pass the right deps to the funcs. We have to build the right dep and pass it to the func. It makes it harder
to use those funcs. We need to make the zero values of those deps more useful. Look at `DemoV6`.
We can check if a dep is nil, then we construct it in the func itself. But we also have DI and if the caller inject the dep,
we can use it.

So now we can still use the func, although we didn't build the deps.

Note: This is how we migrate from the code that doesn't use DI, to code that does use it.

The problem is, checking for deps being nil and construct them everywhere, could become tedious.

Fix: Create a method to construct it. Look at `logger()` method of `ThingV2`. So now all of the funcs or methods that would use that dep,
no longer need to check if the dep is set or not, they can call a method that constructs that dep if not set, or would just return it.

Note:
1. We need to look out for race conditions as well.
2. do we need to construct or set to a default value the new dep everytime after checking it being nil? No, we can use `sync.Once`.

**Note:** When using sync.Once, all the methods should get pointer receivers.

`sync.Once.Do()` is safe for multiple goroutines to call it. Only one of the goroutines will call it and others will wait.

Our solution works except in one scenario which is if we set it and then later somebody set it to nil, this code won't catch that.
To catch that correctly, you'd probably need to add sync.Mutex and then you use that to control access to the dep(like Logger in examples),
to make sure that it's set. But that's overkill.

Note: Making the zero values work especially in a dep of type struct, is a lot of work and it's also not fail-proof. Because
we have to make sure we call the method that returns the constructed dep if it's not set or return the dep if it's actually set. If we forget
to call that method or func in some places and we directly access the field of the dep, we would have problem.

But still adding some more code to make the zero values of deps useful, is worthwhile.

We could also make some of the fields of the dep to be private, so now the only time we might want to replace it, is in a test, using
an internal test.

064 Removing global state with DI

065 Package level functions

066 Summary of dependency injection
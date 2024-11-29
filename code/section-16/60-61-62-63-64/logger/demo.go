package logger

import (
	"errors"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// func OtherWork() {
// 	logger = log.New(os.Stderr, "", 0)
// }

func DemoGlobal() {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

/* DemoV1 has a dep on logger. If we wanna use DI, the way we change this code, is in DemoV2.*/
func DemoV1() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

/*
	Instead of creating the dependency(logger) inside this func, we want the caller to provide it. So this func is not gonna

worry about constructing the dependency on it's own, we want it to be injected(provided) so that it can use it.

So DemoV2 is having logger(dep) injected into it.

Still there's a problem with DemoV2. It requires the caller to pass an EXACT struct(logger) - it's very strict in what it accepts.
*/
func DemoV2(logger *log.Logger) {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

// Call this with:
//
//	logger := log.New(...)
//	DemoV3(log.Println)
//
// Calling this func is tedious, that's why there are better versions of this func later
func DemoV3(logFn func(...interface{})) {
	err := doTheThing()
	if err != nil {
		logFn("error in doTheThing():", err)
	}
}

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

// Call this with:
//
//	logger := log.New(...)
//	DemoV4(logger)
func DemoV4(logger Logger) {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
		logger.Printf("error: %s\n", err)
	}
}

// var defaultThing Thing

// func DemoV5() {
// 	defaultThing.DemoV5()
// }

/*
	Note: With this type, we removed global state. The caller only constructed Thing once and then use it, without any global state.

But we still have state here in this struct in a sense that the fields of this struct could be altered by any methods of this type,
but the upside is we can only change this type is if we pass Thing struct to it. Other code can't change it unlike the global states which
any code can change it.

Note: By using a struct type, it makes it easier to set everything once and then call a bunch of methods without having to
pass in the fields everytime(inject the deps in every single method call), but the cost is giving up a bit of clarity. Since
the caller doesn't know which fields are gonna be used in which methods of this struct.
*/
type Thing struct {
	Logger interface {
		Println(...interface{})
		Printf(string, ...interface{})
	}
	// Printer interface {
	// 	...
	// }
}

func (t Thing) DemoV5() {
	err := doTheThing()
	if err != nil {
		t.Logger.Println("error in doTheThing():", err)
		t.Logger.Printf("error: %s\n", err)
	}
}

func DemoV6(logger Logger) {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
		logger.Printf("error: %s\n", err)
	}
}

type ThingV2 struct {
	Logger interface {
		Println(...interface{})
		Printf(string, ...interface{})
	}
	once sync.Once
}

func (t *ThingV2) logger() Logger {
	t.once.Do(func() {
		if t.Logger == nil {
			t.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		}
	})
	return t.Logger
}

func (t *ThingV2) DemoV7() {
	err := doTheThing()
	if err != nil {
		t.logger().Println("error in doTheThing():", err)
		t.logger().Printf("error: %s\n", err)
	}
}

func doTheThing() error {
	// return nil
	return errors.New("error opening file: abc.txt")
}

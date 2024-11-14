package _25

import (
	"fmt"
	"testing"
	"time"
)

// teardowns will be called before tests are finished.
func TestB(t *testing.T) {
	fmt.Println("setup")
	defer fmt.Println("deferred teardown")

	t.Run("sub1", func(t *testing.T) {
		t.Parallel()
		time.Sleep(time.Second)
		fmt.Println("sub1 done")
	})

	t.Run("sub2", func(t *testing.T) {
		t.Parallel()
		time.Sleep(time.Second)
		fmt.Println("sub2 done")
	})

	fmt.Println("teardown")
}

// solution:
func TestB2(t *testing.T) {
	fmt.Println("setup")
	defer fmt.Println("deferred teardown")

	t.Run("group", func(t *testing.T) {
		t.Run("sub1", func(t *testing.T) {
			t.Parallel()
			fmt.Println("DOING 1")
			time.Sleep(time.Second)
			fmt.Println("sub1 done")
		})

		t.Run("sub2", func(t *testing.T) {
			t.Parallel()
			fmt.Println("DOING 2")
			time.Sleep(time.Second)
			fmt.Println("sub2 done")
		})
	})

	fmt.Println("teardown")
}

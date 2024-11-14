package _24

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
}

func TestA(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
}

func TestB(t *testing.T) {
	t.Run("sub1", func(t *testing.T) {
		t.Parallel()
		// run sub1
	})

	t.Run("sub2", func(t *testing.T) {
		t.Parallel()
		// run sub2
	})
}

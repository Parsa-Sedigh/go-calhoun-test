package _26

import (
	"fmt"
	"testing"
)

// ❌
func TestGotcha(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("i=%d", i), func(t *testing.T) {
			t.Parallel()
			t.Logf("Testing with i=%d", i)
		})
	}
}

// ✅ using local variable for each iteration
func TestGotcha2(t *testing.T) {
	for i := 0; i < 10; i++ {
		/* shadowing, but it's fine here. It's better to use another name for the local var to avoid confusion. */
		// copy value for parallel tests - do not delete this!
		i := i

		t.Run(fmt.Sprintf("i=%d", i), func(t *testing.T) {
			t.Parallel()
			t.Logf("Testing with i=%d", i)
		})
	}
}

type testcase struct {
	arg  int
	want int
}

// ✅ wrapping the closure with another closure
func TestGotcha3(t *testing.T) {
	testCases := []testcase{
		{2, 5}, // clearly the test should fail. But if the fixes are not used, the test will pass since it will use the last tc!
		{3, 9},
		{4, 16},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("arg=%d", tc.arg), closure(tc))
	}
}

func closure(tc testcase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Logf("Testing with: arg=%d, want=%d", tc.arg, tc.want)

		if tc.arg*tc.arg != tc.want {
			t.Errorf("%d^2 != %d", tc.arg, tc.want)
		}
	}
}

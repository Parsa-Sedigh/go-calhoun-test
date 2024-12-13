package stub_test

import (
	"testing"

	"github.com/joncalhoun/twg/suite"
	"github.com/joncalhoun/twg/suite/stub"
	"github.com/joncalhoun/twg/suite/suitetest"
)

var _ suite.UserStore = &stub.UserStore{}

// approach 1 for handling setup and teardown to interface test suite. We have to pass nil
func TestUserStore(t *testing.T) {
	us := &stub.UserStore{}

	// As a client of the interface test suite, we call the test suite func with our own version of the interface
	// that we implemented(us).
	suitetest.UserStore(t, us, nil, nil)
}

// approach 2. Use a struct and expose an All() method instead of calling the interface test suite itself(it should be private in this case)
func TestUserStore_withStruct(t *testing.T) {
	us := &stub.UserStore{}
	tests := suitetest.UserStoreSuite{
		UserStore: us,
	}

	tests.All(t)
}

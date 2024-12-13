package suitetest

import (
	"testing"

	"github.com/joncalhoun/twg/suite"
)

type UserStoreSuite struct {
	suite.UserStore

	BeforeEach func()
	AfterEach  func()
}

func (uss *UserStoreSuite) All(t *testing.T) {
	UserStore(t, uss.UserStore, uss.BeforeEach, uss.AfterEach)
}

// UserStore is an interface test suite - it allows us to test a UserStore implementation which is an interface
// note: we could make this func private, since we expose a type and the All() method.

/* note: We didn't name this func TestUserStore, although it's testing UserStore. That's because the name of this package already
has test in it(suitetest). But if you put this func in the same package as the type(UserStore), we would name it TestUserStore.*/
// suitetest.UserStore
func UserStore(t *testing.T, us suite.UserStore, beforeEach, afterEach func()) {
	_, err := us.ByID(123)
	if err != suite.ErrNotFound {
		t.Errorf("ByID(123) err = nil; want ErrNotFound")
	}

	t.Run("create", func(t *testing.T) {
		user := &suite.User{
			Email: "jon@calhoun.io",
		}
		err = us.Create(user)
		if err != nil {
			t.Errorf("Create() err = %s; want nil", err)
		}
		if user.ID <= 0 {
			t.Errorf("Create() user.ID = %d; want a positive value", user.ID)
		}
	})

	// t.Run("ByID", func(t *testing.T) {
	// 	if beforeEach != nil {
	// 		beforeEach()
	// 	}

	// 	// setup
	// 	user := &suite.User{
	// 		Email: "jon@calhoun.io",
	// 	}

	//////////////////

	// it's not good to test for Create not working correctly in this specific test. Because we should be testing for that scenario
	// in another separate test case and do not concern ourselves with testing that scenario here. Therefore, we don't test if err is not nil here:
	// err = us.Create(user)
	//if err != nil {
	//	t.Errorf("Create() err = %s; want nil", err)
	//}
	//
	//if user.ID <= 0 {
	//	t.Errorf("Create() user.ID = %d; want a positive value", user.ID)
	//}

	//////////////////

	// 	// teardown
	// 	defer func() {
	// 		us.Delete(user)
	// 		if afterEach != nil {
	// 			afterEach()
	// 		}
	// 	}()

	// the actual test cases are after setup and teardown:
	// 	got, err := us.ByID(user.ID)
	// 	if err != nil {
	// 		t.Errorf("ByID() err = %s; want nil", err)
	// 	}
	// 	if got != user {
	// 		t.Errorf("ByID() = %v; want %v", got, user)
	// 	}
	// })

	//
}

package timing

import (
	"testing"
	"time"
)

// We wanna test when we save a user, UpdatedAt is set to the current time.

func TestSaveUser(t *testing.T) {
	now := time.Now()

	/* Here, we're overwriting the global timeNow variable. With this, we're always returning the same `now` variable everytime
	this func is called. So the returned val will be the same. So we will have deterministic test.*/
	timeNow = func() time.Time {
		return now
	}
	defer func() {
		timeNow = time.Now
	}()

	user := User{}
	SaveUser(&user)
	if user.UpdatedAt != now {
		t.Errorf("user.UpdatedAt = %v, want ~%v", user.UpdatedAt, now)
	}
}

func TestUserSaver_Save(t *testing.T) {
	now := time.Now()
	us := UserSaver{
		now: func() time.Time {
			return now
		},
	}
	user := User{}
	us.Save(&user)
	if user.UpdatedAt != now {
		t.Errorf("user.UpdatedAt = %v, want ~%v", user.UpdatedAt, now)
	}
}

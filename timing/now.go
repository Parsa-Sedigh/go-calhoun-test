package timing

import (
	"time"
)

var (
	timeNow = time.Now
)

type User struct {
	UpdatedAt time.Time
}

func SaveUser(user *User) {
	t := timeNow()
	user.UpdatedAt = t
	// ... save the user
}

type UserSaver struct {
	/* Since here we have a struct, we don't need a global time variable like timeNow var. We can just put the func for returning
	a time as a field of this struct. So later we can mock it in our tests.
	We also don't need to make this field public.*/
	now func() time.Time
}

func (us *UserSaver) Save(user *User) {
	t := us.now()
	user.UpdatedAt = t
	// ... save the user
}

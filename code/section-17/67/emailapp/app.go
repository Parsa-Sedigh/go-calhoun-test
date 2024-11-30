package emailapp

import "strings"

// Or SendGridClient - I just happen to use Mailgun.
type MailgunClient struct {
	// stuff here
}

func (mc *MailgunClient) Welcome(name, email string) error {
	// send out a welcome email to the user!
	return nil
}

// this is all fake just to make the demo work
type User struct{}
type UserStore struct{}

func (us *UserStore) Create(name, email string) (*User, error) {
	// pretend to add user to DB
	return &User{}, nil
}

type EmailClient interface {
	Welcome(name, email string) error
}

/*
	Do not use a concrete type for the email client, like MailgunClient. Instead, use an interface like EmailClient.

If you use concrete types(which would have concrete impl of their methods), that means in tests, there wouldn't be a good way to
test without actually sending out emails.

Note: If you replace MailgunClient param(dep) with another real implementation, it's not considered mocking at all,
it's just dep injection and you're just swapping implementations.
*/
func Signup(name, email string, ec EmailClient, us *UserStore) (*User, error) {
	email = strings.ToLower(email)
	user, err := us.Create(name, email)
	if err != nil {
		return nil, err
	}
	err = ec.Welcome(name, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

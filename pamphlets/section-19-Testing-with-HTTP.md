## Section 19 Testing with HTTP
## 075 httptestResponseRecorder
How to test http servers.

Most of the section will cover the <https://golang.org/pkg/net/http/httptest> package, but will also discuss other ideas and suggestions.

### httptest.ResponseRecorder
<https://golang.org/pkg/net/http/httptest#ResponseRecorder>

> ResponseRecorder is an implementation of http.ResponseWriter that records it's mutations for later inspection in tests.

Look at `app_test.go` `TestHome` func.

Testing the response body as a string and comparing it with what we want, in an exact way is not great! Look at TestHome func.
A better way for testing this, is to use parsers(html parsers) to check for certain tags and ... . For example, verify we get a 
specific html tag.

Whenever you wanna test a handler directly, you usually gonna use `httptest.ResponseRecorder`. and `httptest.NewRequest()` and pass these
to the handler func and call it in the test, get the response, read it and see the body or status code or ... .

## 076 httptestServer
The reasons we want to use a test server are:
1. we wanna test our app with a test server - using a request, we 
2. the handler we wanna test is not exported. So we can't call it in our test(since it's in x_test package). So we can't construct
`httptest.ResponseRecorder` and `httptest.NewRecorder`. Note that we can construct the type that has that private handler, but we don't want
to do that. We want to create a test server. For example, in `app` dir > app_test.go, we could construct the Server type and call the
home() method on it, but we don't want to do that. We want to spin up a test web server. We wanna make a real req to it.

Q: When to use which? So when to call the handler directly(without using or spinning up test server) vs spinning up a test server?

A: If you wanna interact directly with the handler, for example if you wanna skip some middleware or other stuff that you might have
set up with server, then calling the handler directly by using httptest.NewRecorder and httptest.NewRequest, is a nice way to just test
specific handler.

But if you wanna test e2e and actually hit the server, use httptest.NewServer(), make an actual req like you're a real user and verify
things happen as expected.

## 077 Building HTTP helpers
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

## 077 Building HTTP helpers
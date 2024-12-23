package app_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/joncalhoun/twg/app"
	"golang.org/x/net/publicsuffix"
)

// call the handler that we wanna test, DIRECTLY without setting up a server
func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.Home(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
	}

	got := string(body)
	want := "<h1>Welcome!</h1>"

	if got != want {
		t.Errorf("GET / = %s; want %s", got, want)
	}
}

// call the handler that we wanna test, using the test server that we have set up
func TestApp_v1(t *testing.T) {
	server := httptest.NewServer(&app.Server{})
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("GET / err = %s; want nil", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
	}

	got := string(body)
	want := "<h1>Welcome!</h1>"

	if got != want {
		t.Errorf("GET / = %s; want %s", got, want)
	}
}

// approach 1: instead of having a "signed in client", we can create a "signed in req" every time we wanna make a req. But this approach has the
// drawback of requiring us specifying a hard-coded cookie session and value, which sometimes we can't have.
func signedInRequest(t *testing.T, method, target string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatalf("http.NewRequest() err = %s; want nil", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "fake_session_token",
	})

	return req
}

// approach 2: instead of using a "signed in req" every time, we can create a "signed in client".
// Note: signedInClient returns a client that has called the login endpoint and have got a session cookie set in it's cookie jar, so that
// we can make future AUTHENTICATED reqs
func signedInClient(t *testing.T, baseURL string) *http.Client {
	// Our cookiejar will keep and set cookies for us between requests. It stores the cookies that the server sets for us and it's gonna
	// set those cookies on every req. So every req that we make with client.Do(), is gonna set those cookies.
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookejar.New() err = %s; want nil", err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// Our client has a cookie jar, but it has no session cookie. By logging
	// in we can ensure that it gets set.
	loginURL := baseURL + "/login"
	req, err := http.NewRequest(http.MethodPost, loginURL, nil)
	if err != nil {
		t.Fatalf("NewRequest() err = %s; want nil", err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("POST /login err = %s; want nil", err)
	}

	t.Logf("Cookies: %v", client.Jar.Cookies(req.URL))

	return client
}

type headerClient struct {
	headers map[string]string
}

/* sets the headers that are in headerClient.headers field, so we don't have to set them everytime in our tests. */
func (hc headerClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for hk, hv := range hc.headers {
		req.Header.Set(hk, hv)
	}

	client := http.Client{}

	return client.Do(req)
}

func TestApp_v2(t *testing.T) {
	server := httptest.NewServer(&app.Server{})
	defer server.Close()

	t.Run("custom built request", func(t *testing.T) {
		t.Log(server.URL)

		req := signedInRequest(t, http.MethodGet, server.URL+"/admin", nil)

		var client http.Client

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("GET /admin err = %s; want nil", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
		}

		got := string(body)
		want := "<h1>Welcome to the admin page!</h1>"

		if got != want {
			t.Errorf("GET /admin = %s; want %s", got, want)
		}
	})

	/* This test is saying: if we go to /admin, we should get 200 since we have the cookie in jar. But it doesn't set the auth header,
	so we should get a 403.*/
	t.Run("cookie based auth", func(t *testing.T) {
		client := signedInClient(t, server.URL)
		res, err := client.Get(server.URL + "/admin")
		if err != nil {
			t.Errorf("GET /admin err = %s; want nil", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("GET /admin code = %d; want %d", res.StatusCode, 200)
		}

		res, err = client.Get(server.URL + "/header-admin")
		if err != nil {
			t.Errorf("GET /header-admin err = %s; want nil", err)
		}
		if res.StatusCode != 403 {
			t.Errorf("GET /header-admin code = %d; want %d", res.StatusCode, 403)
		}
	})

	t.Run("header based auth", func(t *testing.T) {
		client := headerClient{
			headers: map[string]string{"api-key": "fake_api_key"},
		}
		res, err := client.Get(server.URL + "/admin")
		if err != nil {
			t.Errorf("GET /admin err = %s; want nil", err)
		}

		if res.StatusCode != 403 {
			t.Errorf("GET /admin code = %d; want %d", res.StatusCode, 403)
		}

		res, err = client.Get(server.URL + "/header-admin")
		if err != nil {
			t.Errorf("GET /header-admin err = %s; want nil", err)
		}

		// we do have the correct auth headers, so we expect a 200, otherwise fail the test.
		if res.StatusCode != 200 {
			t.Errorf("GET /header-admin code = %d; want %d", res.StatusCode, 200)
		}
	})
}

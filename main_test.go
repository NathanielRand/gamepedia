package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// Form a new HTTP request. This is the request that is going to be
	// passed into our handler. The first arugment is the method, the second
	// argument is the route, and the third argument is the request body.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We will use Go's httptest library to create an http recorder.
	// This recorder will act as the target of our http request.
	// (Think of this as a mini-browser, which will accept the result of the
	// http request that we make.)
	recorder := httptest.NewRecorder()

	// Create an HTTP handler from our handler function. "handler" is the
	// handler function defined in our main.go that we want to test.
	hf := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. THis is the line that actually
	// executes our handler that we want to test.
	hf.ServeHTTP(recorder, req)

	// Check that the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned the wrong status code: got %v when we want %v", status, http.StatusOK)
	}

	// Check that the response body is what we expect.
	expected := `Hello Gamepedia!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v when we want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	// Instantiate the router using the constructor function
	// defined in the main.go file.
	r := newRouter()

	// Create a mock server using the "httptest" libraries `NewServer` method.
	// Documentation: https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// The mock server we created runs a sever and exposes its location
	// in the URL attribute. We make a GET request to the "hello" route
	// we defined in the router.
	resp, err := http.Get(mockServer.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (OK).
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	// In the next few lines, the response body is read, and converted to
	// string.
	defer resp.Body.Close()

	// Read the body into a bunch of bytes (b)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the bytes to a string.
	respString := string(b)
	expected := "Hello Gamepedia!"

	// We want our response to match the one defined in our handler.
	// If it does happen to be "Hello Gamepedia!", then it confirms that the
	// route is correct.
	if respString != expected {
		t.Errorf("Response should be %s, but we go %s", expected, respString)
	}
}

func TestRouterForNonExistantRoute(t *testing.T) {
	// Instantiate the router using the constructor function
	// defined in the main.go file.
	r := newRouter()

	// Create a mock server using the "httptest" libraries `NewServer` method.
	// Documentation: https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// Make a request to a route we know we didn't define, like
	// `POST /hello` route.
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 405 (METHOD NOT ALLOWED)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Response should be 405, but we got %d instead", resp.StatusCode)
	}

	// Defer close of the response body until this function returns.
	defer resp.Body.Close()

	// Read the body into a bunch of bytes (b)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the bytes to a string.
	respString := string(b)
	expected := ""

	// We expect the response to return an empty body.
	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

func TestStaticFileServer(t *testing.T) {
	// Instantiate the router using the constructor function
	// defined in the main.go file.
	r := newRouter()

	// Create a mock server using the "httptest" libraries `NewServer` method.
	// Documentation: https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// We want to hit the "GET /assets/" route to get the index.html file
	// response.
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200 OK, got %d instead", resp.StatusCode)
	}

	// It isn't wise to test the entire conten of the HTML file.
	// Instead, we test that the content-type header is
	// "text/html; charset=utf-8" so that we know that an
	// html file has been served.
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}

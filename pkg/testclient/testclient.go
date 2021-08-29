package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gorilla/schema"
)

// TestClient provides a wrapper for communication with endpoints.
type TestClient struct {
	http *http.Client
	// router  *echo.Echo
	// router  Router
	router  http.Handler
	encoder *schema.Encoder
	token   string
}

// Router provides behavior interface for router object.
// type Router interface {
// 	ServeHTTP(w http.ResponseWriter, r *http.Request)
// }

// Options is used to parameterize new TestClient instance.
type Options struct {
	// Router *echo.Echo
	// Router Router
	Router http.Handler
	Token  string
}

// Setup is used to initialize TestClient with provided options.
func (t *TestClient) Setup(options *Options) {
	encoder := schema.NewEncoder()
	encoder.SetAliasTag("form")

	t.router = options.Router
	t.token = options.Token
	t.encoder = encoder
	t.http = &http.Client{
		Timeout: time.Second * 10,
	}
}

// RequestOptions is used to parameterize Request.
type RequestOptions struct {
	Method   string
	URL      string
	Query    interface{}
	Body     interface{}
	Headers  map[string]string
	Response interface{}
}

// Request provides generic method for http sending requests to API and parsing response body.
func (t *TestClient) Request(options *RequestOptions) error {
	recorder := httptest.NewRecorder()

	// encode query
	if options.Query != nil {
		query := url.Values{}
		err := t.encoder.Encode(options.Query, query)
		if err != nil {
			return fmt.Errorf("failed to encode query: %s", err)
		}

		// append query to url
		options.URL = fmt.Sprintf("%s?%s", options.URL, query.Encode())
	}

	// encode body
	if options.Body != nil {
		bodyString, err := json.Marshal(options.Body)
		if err != nil {
			return fmt.Errorf("failed to encode body, %s", err)
		}

		// rewrite body as a slice of bytes
		options.Body = bytes.NewReader(bodyString)
	} else {
		// or set empty if not provided
		options.Body = bytes.NewReader([]byte(""))
	}

	// create request, add query to url
	request, err := http.NewRequest(
		options.Method,
		options.URL,
		options.Body.(io.Reader),
	)
	if err != nil {
		return fmt.Errorf("failed to send request, %s", err)
	}

	// set request encoding type
	request.Header.Add("Content-type", "application/json")
	// if token provided set auth header
	if t.token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}

	// set custom headers
	for key, value := range options.Headers {
		request.Header.Add(key, value)
	}

	// send request and record response
	t.router.ServeHTTP(recorder, request)

	// convert body to string
	bodyString := recorder.Body.String()

	// if ok
	if recorder.Code >= 200 && recorder.Code < 300 {
		// unmarshal result string into output
		err = json.Unmarshal([]byte(bodyString), &options.Response)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body, %s, %s", bodyString, err)
		}

		return nil
	}

	return fmt.Errorf("received error from server: %s (%d)", bodyString, recorder.Code)
}

package testclient

import (
	"github.com/jarcoal/httpmock"
)

// stub libs
// https://github.com/jarcoal/httpmock
// https://github.com/go-resty/resty

// StubHandler represents a handler stub.
type StubHandler struct {
	Method  string
	URL     string
	Handler httpmock.Responder
}

// StubService represents a stub service.
type StubService struct {
	Name     string
	Handlers []StubHandler
}

// StubServices is used to stub services with test responses.
func StubServices(services []StubService) func() {
	// setup mock server
	httpmock.Activate()

	// register stub handlers
	for _, s := range services {
		for _, h := range s.Handlers {
			httpmock.RegisterResponder(
				h.Method,
				h.URL,
				h.Handler,
			)
		}
	}

	// return teardown for mock server
	return func() {
		httpmock.DeactivateAndReset()
	}
}

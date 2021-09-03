package tests

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/jarcoal/httpmock"
)

// StubServices provides stubbing to third party services.
func StubServices() func() {

	googleFixtures := GetGoogleFixtures()
	facebookFixtures := GetFacebookFixtures()
	githubFixtures := GetGithubFixtures()

	return testclient.StubServices([]testclient.StubService{
		{
			Name: "google-oauth2",
			Handlers: []testclient.StubHandler{
				{
					Method: http.MethodPost,
					URL:    "https://oauth2.googleapis.com/token",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &googleFixtures.Token)
					},
				},
				{
					Method: http.MethodGet,
					URL:    "https://www.googleapis.com/oauth2/v2/userinfo",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &googleFixtures.GoogleUser)
					},
				},
			},
		},

		{
			Name: "facebook-oauth2",
			Handlers: []testclient.StubHandler{
				{
					Method: http.MethodPost,
					URL:    "https://graph.facebook.com/v3.2/oauth/access_token",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &facebookFixtures.Token)
					},
				},
				{
					Method: http.MethodGet,
					URL:    "https://graph.facebook.com/me",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &facebookFixtures.FacebookUser)
					},
				},
			},
		},
		{
			Name: "github-oauth2",
			Handlers: []testclient.StubHandler{
				{
					Method: http.MethodPost,
					URL:    "https://github.com/login/oauth/access_token",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &githubFixtures.Token)
					},
				},
				{
					Method: http.MethodGet,
					URL:    "https://api.github.com/user",
					Handler: func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, &githubFixtures.GithubUser)
					},
				},
			},
		},
	})
}

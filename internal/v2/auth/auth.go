package auth

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	*service.Service
	GoogleOAuthConfig   *oauth2.Config
	FacebookOauthConfig *oauth2.Config
	GithubOauthConfig   *oauth2.Config
}

func (a *Auth) Setup(s *service.Service) {
	a.Service = s

	// setup google oauth config
	a.GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  a.Config.GoogleRedirectURL,
		ClientID:     a.Config.GoogleClientID,
		ClientSecret: a.Config.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// configure router
	GoogleAuthRouter := a.Echo.Group("/api/v2/auth/google")
	GoogleAuthRouter.GET("/login", a.GoogleLoginHandler)
	GoogleAuthRouter.GET("/callback", a.GoogleCallbackHandler)

	// setup facebook oauth config
	a.FacebookOauthConfig = &oauth2.Config{
		ClientID:     a.Config.FacebookClientID,
		ClientSecret: a.Config.FacebookClientSecret,
		RedirectURL:  a.Config.FacebookRedirectURL,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}

	// configure router
	FacebookAuthRouter := a.Echo.Group("/api/v2/auth/facebook")
	FacebookAuthRouter.GET("/login", a.FacebookLoginHandler)
	FacebookAuthRouter.GET("/callback", a.FacebookCallbackHandler)

	// setup github oauth config
	a.GithubOauthConfig = &oauth2.Config{
		ClientID:     a.Config.GithubClientID,
		ClientSecret: a.Config.GithubClientSecret,
		RedirectURL:  a.Config.GithubRedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// configure router
	GithubAuthRouter := a.Echo.Group("/api/v2/auth/github")
	GithubAuthRouter.GET("/login", a.GithubLoginHandler)
	GithubAuthRouter.GET("/callback", a.GithubCallbackHandler)
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
func (p *Auth) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
	// check accept header
	accept := c.Request().Header["Accept"]
	if len(accept) == 0 {
		// default response if accept header is not provided
		return c.JSON(statusCode, res)
	}

	// based on first value in accept header write response
	switch accept[0] {
	case string(models.MimeTypesXML):
		// response with xml
		return c.XML(statusCode, res)
	default:
		// default response with json
		return c.JSON(statusCode, res)
	}
}

package auth

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	*service.Service
	GoogleOAuthConfig *oauth2.Config
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

	// FacebookAuthRouter := a.Echo.Group("/api/v2/auth/facebook")
	// FacebookAuthRouter.GET("/login", a.FacebookLoginHandler)
	// FacebookAuthRouter.GET("/callback", a.FacebookCallbackHandler)

	// TwitterAuthRouter := a.Echo.Group("/api/v2/auth/twitter")
	// TwitterAuthRouter.GET("/login", a.TwitterLoginHandler)
	// TwitterAuthRouter.GET("/callback", a.TwitterCallbackHandler)
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

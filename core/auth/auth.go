package auth

import (
	"fmt"

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

func (a *Auth) Setup(service *service.Service) {
	a.Service = service

	// setup auth config
	a.GoogleOAuthConfig = &oauth2.Config{
		RedirectURL: fmt.Sprintf("%s:8000/api/v2/auth/google/callback", a.Config.BaseURL /*, a.Config.Port*/),
		// ClientID:     m.Config.GoogleClientID,
		// ClientSecret: m.Config.GoogleClientSecret,
		ClientID:     a.Config.GoogleClientID,
		ClientSecret: a.Config.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// configure router
	AuthRouter := a.Echo.Group("/api/v2/auth/google")

	AuthRouter.GET("/login", a.GoogleLoginHandler)
	AuthRouter.GET("/callback", a.GoogleCallbackHandler)
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
func (p *Auth) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
	// check accept header
	accept := c.Request().Header["Accept"][0]

	// based on first value in accept header write response
	switch accept {
	case string(models.MimeTypesXML):
		// response with xml
		return c.XML(statusCode, res)
	default:
		// default response with json
		return c.JSON(statusCode, res)
	}
}

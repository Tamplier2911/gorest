package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// GoogleLoginHandler godoc
//
// @id				GoogleLogin
// @Summary 		Login with google.
// @Description 	Directs users to google popup to grant access to user account.
//
// @Tags			Auth
//
// @Success 307 	{string} Url "url"
//
// @Router /auth/google/login [GET]
func (a *Auth) GoogleLoginHandler(c echo.Context) error {
	// get authorization grant
	logger := a.Logger.Named("GoogleLoginHandler")
	url := a.GoogleOAuthConfig.AuthCodeURL(a.Config.GoogleClientState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	logger = logger.With("url", url)

	logger.Infow("successfully created redirect url")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

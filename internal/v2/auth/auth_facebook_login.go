package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// FacebookLoginHandler godoc
//
// @id				FacebookLogin
// @Summary 		Login with facebook.
// @Description 	Directs users to facebook popup to grant access to user account.
//
// @Tags			Auth
//
// @Success 307 	{string} Url "url"
//
// @Router /auth/facebook/login [GET]
func (a *Auth) FacebookLoginHandler(c echo.Context) error {
	logger := a.Logger.Named("FacebookLoginHandler")

	// force dialog window
	ForceDialog := oauth2.SetAuthURLParam("auth_type", "rerequest")
	// reauthorize - always has for permissions
	// rerequest - for declined/revoked permissions
	// reauthenticate - always as user to confirm password

	// get authorization grant
	url := a.FacebookOauthConfig.AuthCodeURL(a.Config.FacebookClientState, ForceDialog)
	logger = logger.With("url", url)

	logger.Infow("successfully created redirect url")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

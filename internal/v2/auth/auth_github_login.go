package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// GithubLoginHandler godoc
//
// @id				GithubLogin
// @Summary 		Login with github.
// @Description 	Directs users to github popup to grant access to user account.
//
// @Tags			Auth
//
// @Success 307 	{string} Url "url"
//
// @Router /auth/github/login [GET]
func (a *Auth) GithubLoginHandler(c echo.Context) error {
	// get authorization grant
	logger := a.Logger.Named("GithubLoginHandler")
	url := a.GithubOauthConfig.AuthCodeURL(a.Config.GithubClientState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	logger = logger.With("url", url)

	logger.Infow("successfully created redirect url")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

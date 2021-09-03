package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represents output body of GithubCallbackHandler
type GithubCallbackHandlerResponseBody struct {
	Token   *string `json:"token" xml:"token"`
	Message string  `json:"message" xml:"message"`
} // @GithubCallbackResponse

// GithubCallbackHandler godoc
//
// @id				GithubCallback
// @Summary 		Callback triggered once user respond to github authorization popup.
// @Description 	Verifies code and state, exchanges code with authorization token,
//					requests resource server for user personal information, stores users personal
//					data in database, signs JWT and responds with signed JWT token.
//
// @Tags			Auth
//
// @Produce json
// @Produce xml
//
// @Param code query string true "Parameter for code grant"
// @Param state query string true "Parameter for state"
//
// @Success 200 	{object} GithubCallbackHandlerResponseBody
// @Failure 400,401 {object} GithubCallbackHandlerResponseBody
// @Failure 500 	{object} GithubCallbackHandlerResponseBody
// @Failure default {object} GithubCallbackHandlerResponseBody
//
// @Router /auth/github/callback [GET]
func (a *Auth) GithubCallbackHandler(c echo.Context) error {
	logger := a.Logger.Named("GithubCallbackHandler")

	// get state from query
	state := c.QueryParam("state")
	logger = logger.With("state", state)

	// get authorization grant from query
	code := c.QueryParam("code")
	logger = logger.With("code", code)

	// check state
	logger.Infow("checking state")
	if state != a.Config.GithubClientState {
		logger.Errorw("invalid state", "state", state)
		return a.ResponseWriter(c, http.StatusUnauthorized, GithubCallbackHandlerResponseBody{
			Message: "invalid auth state",
		})
	}

	// exchange authorization grant with token
	logger.Infow("exchanging token")
	token, err := a.GithubOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		logger.Errorw("failed to exchange token", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GithubCallbackHandlerResponseBody{
			Message: fmt.Sprintf("code exchange failed: %s", err.Error()),
		})
	}
	logger = logger.With("token", token)

	// create client request for github user data using authorization token
	oauthClient := a.GithubOauthConfig.Client(context.TODO(), token)
	client := github.NewClient(oauthClient)
	ghu, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		logger.Errorw("failed to get user info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GithubCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed getting user info: %s", err.Error()),
		})
	}
	logger = logger.With("githubUser", ghu)

	if *ghu.Email == "" {
		logger.Errorw("user does not have email address")
		return a.ResponseWriter(c, http.StatusForbidden, GithubCallbackHandlerResponseBody{
			Message: "email address is required, make sure you have public email address set in your github account",
		})
	}

	logger.Infow("successfully authorized with github")

	// get user from database
	logger.Infow("getting user from database")
	var user models.User
	err = a.MySQL.
		Model(&models.User{}).
		Where(&models.User{Email: *ghu.Email}).
		First(&user).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find user in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if user not found create new user record
	if err == gorm.ErrRecordNotFound {
		logger.Infow("could not find user with this email, creating new user record")
		// create user record in database
		user = models.User{
			Email:     *ghu.Email,
			Username:  *ghu.Name,
			AvatarURL: *ghu.AvatarURL,
			UserRole:  models.UserRoleUser,
		}
		err := a.MySQL.Create(&user).Error
		if err != nil {
			logger.Errorw("failed to create new user in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
				Message: "failed to register new user",
			})
		}
	}
	logger = logger.With("user", user)

	// get auth provider from database
	logger.Infow("getting auth provider from database")
	var authProvider models.AuthProvider
	err = a.MySQL.
		Model(&models.AuthProvider{}).
		Where(&models.AuthProvider{
			UserID:           user.ID,
			ProviderUID:      fmt.Sprintf("%d", *ghu.ID),
			AuthProviderType: models.AuthProviderTypeGithub,
		}).
		First(&authProvider).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find auth provider in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if auth provider token found update short living refresh token
	if err == nil {
		// TODO: consider encrypt token before updating
		// ensure to update short living token in database
		logger.Infow("updating auth provider in database")
		err = a.MySQL.
			Model(&authProvider).
			Updates(&models.AuthProvider{RefreshToken: token.AccessToken}).
			Error
		if err != nil {
			logger.Errorw("failed to update auth provider in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
				Message: "failed to login user",
			})
		}
	}

	// if auth provider not found create new auth provider
	if err == gorm.ErrRecordNotFound {
		// TODO: consider encrypt token before saving
		if token.AccessToken != "" {
			logger.Infow("saving auth provder to database")
			authProvider = models.AuthProvider{
				UserID:           user.ID,
				ProviderUID:      fmt.Sprintf("%d", *ghu.ID),
				AuthProviderType: models.AuthProviderTypeGithub,
				RefreshToken:     token.AccessToken,
			}
			err := a.MySQL.Create(&authProvider).Error
			if err != nil {
				logger.Errorw("failed to create new auth provider in database", "err", err)
				return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
					Message: "failed to login user",
				})
			}
		}
	}
	logger = logger.With("authProvider", authProvider)

	// sign jwt token
	logger.Infow("encoding jwt token")
	jwt, err := access.EncodeToken(&access.Token{
		UserID:   user.ID,
		UserRole: user.UserRole,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "gorest-api",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
		},
	}, a.Config.HMACSecret)
	if err != nil {
		logger.Errorw("failed to sign jwt token", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, GithubCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}
	logger = logger.With("jwt token", jwt)

	logger.Infow("successfully logged in")
	return a.ResponseWriter(c, http.StatusOK, GithubCallbackHandlerResponseBody{
		Token:   &jwt,
		Message: "successfully logged in",
	})
}

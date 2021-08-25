package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represents output body of FacebookCallbackHandler
type FacebookCallbackHandlerResponseBody struct {
	Token   *string `json:"token" xml:"token"`
	Message string  `json:"message" xml:"message"`
} // @FacebookCallbackResponse

// Represents user object returned from facebook resource server
type FacebookUserData struct {
	ID    string `json:"id" xml:"id"`
	Email string `json:"email" xml:"email"`
	Name  string `json:"name" xml:"name"`
} // @name FacebookUserData

// FacebookCallbackHandler godoc
//
// @id				FacebookCallback
// @Summary 		Callback triggered once user respond to facebook authorization popup.
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
// @Success 200 	{object} FacebookCallbackHandlerResponseBody
// @Failure 400,401 {object} FacebookCallbackHandlerResponseBody
// @Failure 500 	{object} FacebookCallbackHandlerResponseBody
// @Failure default {object} FacebookCallbackHandlerResponseBody
//
// @Router /auth/facebook/callback [GET]
func (a *Auth) FacebookCallbackHandler(c echo.Context) error {

	logger := a.Logger.Named("FacebookCallbackHandler")

	// get state from query
	state := c.QueryParam("state")
	logger = logger.With("state", state)

	// get authorization grant from query
	code := c.QueryParam("code")
	logger = logger.With("code", code)

	// check state
	logger.Infow("checking state")
	if state != a.Config.FacebookClientState {
		logger.Errorw("invalid state", "state", state)
		return a.ResponseWriter(c, http.StatusUnauthorized, FacebookCallbackHandlerResponseBody{
			Message: "invalid auth state",
		})
	}

	// exchange authorization grant with token
	logger.Infow("exchanging token")
	token, err := a.FacebookOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		logger.Errorw("failed to exchange token", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, FacebookCallbackHandlerResponseBody{
			Message: fmt.Sprintf("code exchange failed: %s", err.Error()),
		})
	}
	logger = logger.With("token", token)

	// access resource server using auth token
	logger.Infow("getting user info")
	res, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		logger.Errorw("failed to get user info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, FacebookCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed getting user info: %s", err.Error()),
		})
	}
	defer res.Body.Close()

	// read resource server response
	logger.Infow("parsing user info")
	rd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorw("failed to parse uer info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, FacebookCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed reading response body: %s", err.Error()),
		})
	}
	logger = logger.With("resourceData", rd)

	// unmarshal facebook user into a struct
	logger.Infow("unmarshal user data")
	var fu FacebookUserData
	err = json.Unmarshal(rd, &fu)
	if err != nil {
		logger.Errorw("failed to parse facebook user info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, FacebookCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed reading response body: %s", err.Error()),
		})
	}
	logger = logger.With("facebook user", fu)

	if fu.Email == "" {
		logger.Errorw("user does not have email address")
		return a.ResponseWriter(c, http.StatusForbidden, FacebookCallbackHandlerResponseBody{
			Message: "email address is required",
		})
	}

	logger.Infow("successfully logged with facebook")

	// get user from database
	logger.Infow("getting user from database")
	var user models.User
	err = a.MySQL.
		Model(&models.User{}).
		Where(&models.User{Email: fu.Email}).
		First(&user).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find user in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if user not found create new user record
	if err == gorm.ErrRecordNotFound {
		logger.Infow("could not find user with this email, creating new user record")
		// create user record in database
		user = models.User{
			Email:    fu.Email,
			Username: fu.Name,
			UserRole: models.UserRoleUser,
		}
		err := a.MySQL.Create(&user).Error
		if err != nil {
			logger.Errorw("failed to create new user in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
				Message: "failed to register new user",
			})
		}
	}
	logger = logger.With("user", user)

	// get auth provider form database
	logger.Infow("getting auth provider from database")
	var authProvider models.AuthProvider
	err = a.MySQL.
		Model(&models.AuthProvider{}).
		Where(&models.AuthProvider{
			UserID:           user.ID,
			ProviderUID:      fu.ID,
			AuthProviderType: models.AuthProviderTypeFacebook,
		}).
		First(&authProvider).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find auth token in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if auth provider found update short living refresh token
	if err == nil {
		// TODO: consider encrypt token before updating
		// ensure to update short living token in database
		logger.Infow("updating auth provider token in database")
		err = a.MySQL.
			Model(&authProvider).
			Updates(&models.AuthProvider{RefreshToken: token.AccessToken}).
			Error
		if err != nil {
			logger.Errorw("failed to update auth provider in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
				Message: "failed to login user",
			})
		}
	}

	// if auth provider not found create new auth provider
	if err == gorm.ErrRecordNotFound {
		// TODO: consider encrypt token before saving
		// facebook does not support refresh token, we use short living token, to get long living token
		if token.AccessToken != "" {
			logger.Infow("saving auth provider to database")
			authProvider = models.AuthProvider{
				UserID:           user.ID,
				ProviderUID:      fu.ID,
				AuthProviderType: models.AuthProviderTypeFacebook,
				RefreshToken:     token.AccessToken,
			}
			err := a.MySQL.Create(&authProvider).Error
			if err != nil {
				logger.Errorw("failed to create new auth provider in database", "err", err)
				return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
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
		return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}
	logger = logger.With("jwt token", jwt)

	logger.Infow("successfully logged in")
	return a.ResponseWriter(c, http.StatusInternalServerError, FacebookCallbackHandlerResponseBody{
		Token:   &jwt,
		Message: "successfully logged in",
	})
}

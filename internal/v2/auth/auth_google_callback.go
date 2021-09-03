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

// Represents output body of GoogleCallbackHandler
type GoogleCallbackHandlerResponseBody struct {
	Token   *string `json:"token" xml:"token"`
	Message string  `json:"message" xml:"message"`
} // @GoogleCallbackResponse

// Represents user object returned from google resource server
type GoogleUserData struct {
	ID            string `json:"id" xml:"id"`
	Email         string `json:"email" xml:"email"`
	EmailVerified bool   `json:"verified_email" xml:"verified_email"`
	Picture       string `json:"picture" xml:"picture"`
} // @name GoogleUserData

// GoogleCallbackHandler godoc
//
// @id				GoogleCallback
// @Summary 		Callback triggered once user respond to google authorization popup.
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
// @Success 200 	{object} GoogleCallbackHandlerResponseBody
// @Failure 400,401 {object} GoogleCallbackHandlerResponseBody
// @Failure 500 	{object} GoogleCallbackHandlerResponseBody
// @Failure default {object} GoogleCallbackHandlerResponseBody
//
// @Router /auth/google/callback [GET]
func (a *Auth) GoogleCallbackHandler(c echo.Context) error {
	logger := a.Logger.Named("GoogleCallbackHandler")

	// get state from query
	state := c.QueryParam("state")
	logger = logger.With("state", state)

	// get authorization grant from query
	code := c.QueryParam("code")
	logger = logger.With("code", code)

	// check state
	logger.Infow("checking state")
	if state != a.Config.GoogleClientState {
		logger.Errorw("invalid state", "state", state)
		return a.ResponseWriter(c, http.StatusUnauthorized, GoogleCallbackHandlerResponseBody{
			Message: "invalid auth state",
		})
	}

	// exchange authorization grant with token
	logger.Infow("exchanging token")
	token, err := a.GoogleOAuthConfig.Exchange(context.TODO(), code)
	if err != nil {
		logger.Errorw("failed to exchange token", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GoogleCallbackHandlerResponseBody{
			Message: fmt.Sprintf("code exchange failed: %s", err.Error()),
		})
	}
	logger = logger.With("token", token)

	// access resource server using auth token
	logger.Infow("getting user info")
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		logger.Errorw("failed to get user info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GoogleCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed getting user info: %s", err.Error()),
		})
	}
	defer res.Body.Close()

	// read resource server response
	logger.Infow("parsing user info")
	rd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorw("failed to parse uer info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GoogleCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed reading response body: %s", err.Error()),
		})
	}
	logger = logger.With("resourceData", rd)

	// unmarshal google user into a struct
	logger.Infow("unmarshal user data")
	var gu GoogleUserData
	err = json.Unmarshal(rd, &gu)
	if err != nil {
		logger.Errorw("failed to parse google user info", "err", err)
		return a.ResponseWriter(c, http.StatusUnauthorized, GoogleCallbackHandlerResponseBody{
			Message: fmt.Sprintf("failed reading response body: %s", err.Error()),
		})
	}
	logger = logger.With("google user", gu)

	if gu.Email == "" {
		logger.Errorw("user does not have email address")
		return a.ResponseWriter(c, http.StatusForbidden, GoogleCallbackHandlerResponseBody{
			Message: "email address is required",
		})
	}

	logger.Infow("successfully authorized with google")

	// get user from database
	logger.Infow("getting user from database")
	var user models.User
	err = a.MySQL.
		Model(&models.User{}).
		Where(&models.User{Email: gu.Email}).
		First(&user).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find user in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if user not found create new user record
	if err == gorm.ErrRecordNotFound {
		logger.Infow("could not find user with this email, creating new user record")
		// create user record in database
		user = models.User{
			Email:     gu.Email,
			AvatarURL: gu.Picture,
			UserRole:  models.UserRoleUser,
		}
		err := a.MySQL.Create(&user).Error
		if err != nil {
			logger.Errorw("failed to create new user in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
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
			ProviderUID:      gu.ID,
			AuthProviderType: models.AuthProviderTypeGoogle,
		}).
		First(&authProvider).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorw("failed to find auth provider in database", "err", err)
		return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}

	// if auth provider not found create new auth provider
	if err == gorm.ErrRecordNotFound {
		// TODO: consider encrypting token before saving
		// save refresh token in order if we want to request resource api when user is offline
		logger.Infow("saving auth provder to database")
		authProvider := models.AuthProvider{
			UserID:           user.ID,
			ProviderUID:      gu.ID,
			AuthProviderType: models.AuthProviderTypeGoogle,
			RefreshToken:     token.RefreshToken,
		}
		err := a.MySQL.Create(&authProvider).Error
		if err != nil {
			logger.Errorw("failed to create new auth provider in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
				Message: "failed to login user",
			})
		}

	}

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
		return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
			Message: "failed to login user",
		})
	}
	logger = logger.With("jwt token", jwt)

	logger.Infow("successfully logged in")
	return a.ResponseWriter(c, http.StatusOK, GoogleCallbackHandlerResponseBody{
		Token:   &jwt,
		Message: "successfully logged in",
	})
}

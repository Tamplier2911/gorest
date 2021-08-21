package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Represents output body of GoogleCallbackHandler
type GoogleCallbackHandlerResponseBody struct {
	User    *GoogleUserData `json:"user" xml:"user"`
	Message string          `json:"message" xml:"message"`
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

	// access resource server using token
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

	logger.Infow("successfully authorized with google")

	// TODO: set all require environment variables and middlewares to pkg in separate PR

	/*
		// get user from database
		logger.Infow("getting user from database")
		var user models.User
		err = a.MySQL.
			Model(&models.User{}).
			Where(&models.User{Email: gu.Email, GoogleUID: gu.ID}).
			First(&user).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Infow("could not find user with this uid, creating new user record")
				// create user record in database
				newUser := models.User{
					Email:     gu.Email,
					AvatarURL: gu.Picture,
					GoogleUID: gu.ID,
				}
				err := a.MySQL.Create(&newUser)
				if err != nil {
					logger.Errorw("failed to create new user in database", "err", err)
					return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
						Message: fmt.Sprintf("failed to register new user"),
					})
				}
				// get user to outer scope
				user = newUser

				// create refresh token record in order we want to request resource api when user is offline
				if token.RefreshToken != "" {
					refreshToken := models.AuthRefreshToken{
						UserID:       newUser.ID,
						AuthProvider: models.AuthProviderGoogle,
						RefreshToken: token.RefreshToken,
					}
					err := a.MySQL.Create(&refreshToken)
					if err != nil {
						logger.Errorw("failed to create new refresh token in database", "err", err)
						return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
							Message: fmt.Sprintf("failed to register new user"),
						})
					}
				}
			}
			logger.Errorw("failed to find user in database", "err", err)
			return a.ResponseWriter(c, http.StatusInternalServerError, GoogleCallbackHandlerResponseBody{
				Message: fmt.Sprintf("failed to get user form database"),
			})
		}

	*/

	// TODO: sign token

	// TODO: respond with a token

	// TODO: once finalized this tasks, implement twitter and facebook oauth as well

	// TODO: be happy

	// return c.JSON(http.StatusOK, CallbackResponseBody{
	// 	User:    &user,
	// 	Message: "success",
	// })
	return c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:8000/api/v2/posts")
}

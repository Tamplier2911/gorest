package auth

import (
	"net/http"
	"time"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represent output data of RefreshTokenHandler
type RefreshTokenHandlerHandlerResponseBody struct {
	Token   *string `json:"token" xml:"token"`
	Message string  `json:"message" xml:"message"`
} // @name RefreshTokenResponse

// RefreshTokenHandler godoc
//
// @id				RefreshToken
// @Summary 		Refreshes token.
// @Description 	Validate user token and produce token with prolonged expire data.
//
// @Tags			Auth
//
// @Produce json
// @Produce xml
//
// @Success 201 		{object} RefreshTokenHandlerHandlerResponseBody
// @Failure 400,403,404 {object} RefreshTokenHandlerHandlerResponseBody
// @Failure 500 		{object} RefreshTokenHandlerHandlerResponseBody
// @Failure default 	{object} RefreshTokenHandlerHandlerResponseBody
//
// @Security ApiKeyAuth
//
// @Router /auth/refresh [POST]
func (a *Auth) RefreshTokenHandler(c echo.Context) error {
	logger := a.Logger.Named("RefreshTokenHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

	// get user from database
	logger.Info("getting user from database")
	var user models.User
	err := a.MySQL.
		Model(&models.User{}).
		Where(&models.User{Base: models.Base{ID: token.UserID}}).
		First(&user).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("failed to find user with id from token", "tokenUserID", token.UserID)
			return a.ResponseWriter(c, http.StatusNotFound, RefreshTokenHandlerHandlerResponseBody{
				Message: "failed to find user",
			})
		}
		logger.Error("failed to find user in database", "tokenUserID", token.UserID)
		return a.ResponseWriter(c, http.StatusInternalServerError, RefreshTokenHandlerHandlerResponseBody{
			Message: "failed to refresh token",
		})
	}
	logger = logger.With("user", user)

	// refresh user token
	logger.Info("refreshing token")
	refreshedToken, err := access.EncodeToken(&access.Token{
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
	logger = logger.With("refreshedToken", refreshedToken)

	// assemble response body
	logger.Infow("assembling response body")
	res := RefreshTokenHandlerHandlerResponseBody{
		Token:   &refreshedToken,
		Message: "successfully refreshed token",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully refreshed token")
	return a.ResponseWriter(c, http.StatusCreated, res)
}

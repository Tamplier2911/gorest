package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthenticationMiddleware is used to authenticate user.
func AuthenticationMiddleware(logger *zap.SugaredLogger, config *config.Config, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := logger.Named("AuthenticationMiddleware")

		// get token and check if header value exist
		authHeaderArr := c.Request().Header["Authorization"]
		logger.Infow("checking authorization header")
		if len(authHeaderArr) == 0 {
			logger.Errorw("empty authorization header", "authHeaderArr", authHeaderArr)
			return echo.NewHTTPError(http.StatusUnauthorized, "empty authorization token")
		}

		// retrieve token from header
		logger.Infow("retrieving token from header")
		tokenArr := strings.Split(authHeaderArr[0], " ")
		if len(tokenArr) != 2 {
			logger.Errorw("malformed auth token", "tokenArr", tokenArr)
			return echo.NewHTTPError(http.StatusUnauthorized, "malformed auth token")
		}

		// get token from header value
		token := tokenArr[1]

		// decode token
		logger.Infow("decoding token", "token", token)
		decodedToken, err := access.DecodeToken(token, config.HMACSecret)
		if err != nil {
			logger.Errorw("failed to decode token", "err", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "malformed auth token")
		}

		// save token to context
		logger.Infow("saving token to context", "decodedToken", decodedToken)
		access.SaveTokenToContext(c, decodedToken)

		// success, pass context to next middleware
		logger.Infow("successfully authenticated request", "decodedToken", decodedToken)
		return next(c)
	}
}

// AuthWrapperDP is used to authenticate user DEPRECATED.
func AuthWrapperDP(
	handler func(w http.ResponseWriter, r *http.Request),
	lg *zap.SugaredLogger,
	cf *config.Config,
	w http.ResponseWriter,
	r *http.Request,
) {
	logger := lg.Named("AuthenticationMiddlewareDP")

	// get token and check if header value exist
	authHeaderStr := r.Header.Get("Authorization")
	logger.Infow("checking authorization header")
	if len(authHeaderStr) == 0 {
		logger.Errorw("empty authorization header", "authHeaderArr", authHeaderStr)
		err := errors.New("empty authorization token")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// retrieve token from header
	logger.Infow("retrieving token from header")
	tokenArr := strings.Split(authHeaderStr, " ")
	if len(tokenArr) != 2 {
		logger.Errorw("malformed auth token", "tokenArr", tokenArr)
		err := errors.New("malformed auth token")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// get token from header value
	token := tokenArr[1]

	// decode token
	logger.Infow("decoding token", "token", token)
	decodedToken, err := access.DecodeToken(token, cf.HMACSecret)
	if err != nil {
		logger.Errorw("failed to decode token", "err", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// save token to context
	logger.Infow("saving token to context", "decodedToken", decodedToken)
	ctx := context.WithValue(r.Context(), "token", decodedToken) //lint:ignore SA1029 deprecated wrapper

	// success, pass context to next middleware
	logger.Infow("successfully authenticated request", "decodedToken", decodedToken)
	handler(w, r.WithContext(ctx))
}

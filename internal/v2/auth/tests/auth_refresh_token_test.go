package tests

import (
	"testing"
	"time"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/auth"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
)

func TestRefreshTokenHandler(t *testing.T) {
	// setup stub
	teardown := StubServices()

	// init service
	a := app.Application{}
	a.Setup()

	// create test user
	user := models.User{
		Username: "refresh auth token",
		Email:    "refresh_auth_token@test.com",
		UserRole: models.UserRoleUser,
	}
	err := a.MySQL.Create(&user).Error
	require.NoError(t, err, "failed to create test user")

	// init test clients
	noTokenClient := testclient.TestClient{}
	noTokenClient.Setup(&testclient.Options{Router: a.Echo})

	noUserInDbClient := testclient.TestClient{}
	noUserInDbClient.Setup(&testclient.Options{
		Router: a.Echo,
		Token:  access.MustEncodeToken(&access.Token{UserID: uuid.New(), UserRole: user.UserRole}, a.Config.HMACSecret),
	})

	yesterday := time.Now().Add(-24 * time.Hour)
	oldToken := access.MustEncodeToken(&access.Token{
		UserID: user.ID, UserRole: user.UserRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: yesterday.Add(time.Hour * 24 * 7).Unix(),
		},
	}, a.Config.HMACSecret)

	authorizedClient := testclient.TestClient{}
	authorizedClient.Setup(&testclient.Options{
		Router: a.Echo,
		Token:  oldToken,
	})

	defer func() {
		// cleanup stubs
		teardown()

		// clean test user
		err := a.MySQL.
			Unscoped().
			Where(&models.User{Base: models.Base{ID: user.ID}}).
			Select(clause.Associations).
			Delete(&models.User{}).
			Error
		require.NoError(t, err, "failed to delete test user")
	}()

	t.Run("should fail if no auth token provided", func(t *testing.T) {
		var res testclient.DefaultResponse
		err := noTokenClient.Request(&testclient.RequestOptions{
			Method:          "GET",
			URL:             "/api/v2/auth/refresh",
			DefaultResponse: &res,
		})
		require.Error(t, err, "unexpected response")
	})

	t.Run("should fail if no user with provided id found in database", func(t *testing.T) {
		var res testclient.DefaultResponse
		err := noUserInDbClient.Request(&testclient.RequestOptions{
			Method:          "GET",
			URL:             "/api/v2/auth/refresh",
			DefaultResponse: &res,
		})
		require.Error(t, err, "unexpected response")
	})

	t.Run("should refresh auth token", func(t *testing.T) {
		var res auth.RefreshTokenHandlerHandlerResponseBody
		err := authorizedClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      "/api/v2/auth/refresh",
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.NotEmpty(t, res.Token, "token empty in response")

		oldToken, _ := access.DecodeToken(oldToken, a.Config.HMACSecret)
		newToken, _ := access.DecodeToken(*res.Token, a.Config.HMACSecret)
		require.Greater(t, newToken.ExpiresAt, oldToken.ExpiresAt, "unexpected expiry date")
		require.NotEmpty(t, newToken.UserID, "token does not have user id")
		require.NotEmpty(t, newToken.UserRole, "token does not have user role")
	})
}

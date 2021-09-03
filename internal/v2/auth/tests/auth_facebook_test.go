package tests

import (
	"fmt"
	"net/http"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/auth"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
)

func TestFacebookAuth(t *testing.T) {
	// setup stub
	teardown := StubServices()

	// init service
	a := app.Application{}
	a.Setup()

	// get facebook fixtures
	facebookFixtures := GetFacebookFixtures()

	// init test client
	testClient := testclient.TestClient{}
	testClient.Setup(&testclient.Options{Router: a.Echo})

	// test entity ids
	var testUserId uuid.UUID

	defer func() {
		// cleanup stubs
		teardown()

		// clean test user
		err := a.MySQL.
			Unscoped().
			Where(&models.User{Base: models.Base{ID: testUserId}}).
			Select(clause.Associations).
			Delete(&models.User{}).
			Error
		require.NoError(t, err, "failed to delete test user")
	}()

	t.Run("should redirect to facebook popup", func(t *testing.T) {
		var res testclient.DefaultResponse
		err := testClient.Request(&testclient.RequestOptions{
			Method:          "GET",
			URL:             "/api/v2/auth/facebook/login",
			DefaultResponse: &res,
		})
		require.NoError(t, err, "parsed invalid uuid")
		require.Equal(t, http.StatusTemporaryRedirect, res.Status, "unexpected response status")
	})

	t.Run("should perform callback logic", func(t *testing.T) {
		var res auth.FacebookCallbackHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/auth/facebook/callback?state=%s", a.Config.FacebookClientState),
			Response: &res,
		})
		require.NoError(t, err, "parsed invalid uuid")
		require.NotEmpty(t, res.Token, "empty token field")
	})

	t.Run("user should be created in database", func(t *testing.T) {
		var user models.User
		err := a.MySQL.
			Model(&models.User{}).
			Where(&models.User{Email: facebookFixtures.FacebookUser.Email}).
			First(&user).
			Error
		require.NoError(t, err, "parsed invalid uuid")
		require.Equal(t, user.UserRole, models.UserRoleUser, "unexpected user role")
		require.Equal(t, user.Username, facebookFixtures.FacebookUser.Name, "unexpected user name")
		require.Equal(t, user.Email, facebookFixtures.FacebookUser.Email, "unexpected email address")
		testUserId = user.ID
	})

	t.Run("auth provider should be created in database", func(t *testing.T) {
		var provider models.AuthProvider
		err := a.MySQL.
			Model(&models.AuthProvider{}).
			Where(&models.AuthProvider{ProviderUID: facebookFixtures.FacebookUser.ID}).
			First(&provider).
			Error
		require.NoError(t, err, "parsed invalid uuid")
		require.Equal(t, provider.AuthProviderType, models.AuthProviderTypeFacebook, "unexpected provider type")
		require.Equal(t, provider.ProviderUID, facebookFixtures.FacebookUser.ID, "unexpected provider uid")
		require.Equal(t, provider.RefreshToken, facebookFixtures.Token.AccessToken, "unexpected refresh token value")
	})

	t.Run("should login user without creating new one", func(t *testing.T) {
		var res auth.FacebookCallbackHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/auth/facebook/callback?state=%s", a.Config.FacebookClientState),
			Response: &res,
		})
		require.NoError(t, err, "parsed invalid uuid")
		require.NotEmpty(t, res.Token, "empty token field")
	})

	t.Run("should be one instance of user in database", func(t *testing.T) {
		var user []models.User
		err := a.MySQL.
			Model(&models.User{}).
			Where(&models.User{Email: facebookFixtures.FacebookUser.Email}).
			Find(&user).
			Error
		require.NoError(t, err, "parsed invalid uuid")
		require.Len(t, user, 1, "unexpected user role")
	})
}

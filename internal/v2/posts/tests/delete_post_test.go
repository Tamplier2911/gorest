package tests

import (
	"fmt"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDeletePostHandler(t *testing.T) {
	// init service
	m := app.Monolith{}
	m.Setup()

	// init test fixtures
	fixture := PostsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	authorClient := testclient.TestClient{}
	authorClient.Setup(&testclient.Options{
		Router: m.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, m.Config.HMACSecret),
	})

	notAuthorClient := testclient.TestClient{}
	notAuthorClient.Setup(&testclient.Options{
		Router: m.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserTwoID,
		}, m.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	t.Run("only post author can delete a post", func(t *testing.T) {
		err := notAuthorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
		})
		require.Error(t, err, "random user can delete post")
	})

	t.Run("should error if passing invalid post id", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v2/posts/%s", uuid.New()),
		})
		require.Error(t, err, "can delete not existing post")
	})

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v2/posts/%s", "invalid uuid"),
		})
		require.Error(t, err, "passing uuid parsing with invalid id")
	})

	t.Run("should delete post", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
		})
		require.NoError(t, err, "failed to delete post")
	})

	t.Run("post should be deleted from", func(t *testing.T) {
		var post models.Post
		err := m.MySQL.
			Model(&models.Post{}).
			Where(&models.Post{Base: models.Base{ID: testData.TestPostOneUserOneID}}).
			First(&post).
			Error
		require.Error(t, err, "got post from database that should be deleted")
		require.Equal(t, gorm.ErrRecordNotFound, err, "found not existing post in database")
	})
}

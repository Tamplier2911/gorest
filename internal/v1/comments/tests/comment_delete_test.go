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

func TestDeleteCommentHandler(t *testing.T) {
	// init service
	a := app.Application{}
	a.Setup()

	// init test fixtures
	fixture := CommentsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	authorClient := testclient.TestClient{}
	authorClient.Setup(&testclient.Options{
		Router: a.Router,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, a.Config.HMACSecret),
	})

	notAuthorClient := testclient.TestClient{}
	notAuthorClient.Setup(&testclient.Options{
		Router: a.Router,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserTwoID,
		}, a.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	t.Run("only comment author can delete a comment", func(t *testing.T) {
		err := notAuthorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
		})
		require.Error(t, err, "random user can delete comment")
	})

	t.Run("should error if passing invalid comment id", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v1/comments/%s", uuid.New()),
		})
		require.Error(t, err, "can delete not existing comment")
	})

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v1/comments/%s", "invalid uuid"),
		})
		require.Error(t, err, "passing uuid parsing with invalid id")
	})

	t.Run("should delete comment", func(t *testing.T) {
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "DELETE",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
		})
		require.NoError(t, err, "failed to delete comment")
	})

	t.Run("comment should be deleted from", func(t *testing.T) {
		var comment models.Comment
		err := a.MySQL.
			Model(&models.Comment{}).
			Where(&models.Comment{Base: models.Base{ID: testData.TestUserOneCommentOneID}}).
			First(&comment).
			Error
		require.Error(t, err, "got comment from database that should be deleted")
		require.Equal(t, gorm.ErrRecordNotFound, err, "found not existing comment in database")
	})
}

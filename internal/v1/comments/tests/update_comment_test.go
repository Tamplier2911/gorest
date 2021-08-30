package tests

import (
	"fmt"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v1/comments"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdateCommentHandler(t *testing.T) {
	// init service
	m := app.Monolith{}
	m.Setup()

	// init test fixtures
	fixture := CommentsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	authorClient := testclient.TestClient{}
	authorClient.Setup(&testclient.Options{
		Router: m.Router,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, m.Config.HMACSecret),
	})

	notAuthorClient := testclient.TestClient{}
	notAuthorClient.Setup(&testclient.Options{
		Router: m.Router,
		Token: access.MustEncodeToken(&access.Token{
			UserID: uuid.New(),
		}, m.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	updatedCommentName := "updated test comment name"
	updatedCommentBody := "updated test comment name"
	t.Run("only author can update comment", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := notAuthorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Name: updatedCommentName,
				Body: updatedCommentBody,
			},
			Response: &res,
		})
		require.Error(t, err, "random user could update comment")
	})

	t.Run("should error if passing invalid post uuid", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", uuid.New()),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Name: updatedCommentName,
				Body: updatedCommentBody,
			},
			Response: &res,
		})
		require.Error(t, err, "updated not existing comment in database")
	})

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", "invalid uuid"),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Name: updatedCommentName,
				Body: updatedCommentBody,
			},
			Response: &res,
		})
		require.Error(t, err, "passing uuid parsing with invalid id")
	})

	t.Run("name field should be required", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Body: updatedCommentBody,
			},
			Response: &res,
		})
		require.Error(t, err, "passing through required check")
	})

	t.Run("body field should be required", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Name: updatedCommentName,
			},
			Response: &res,
		})
		require.Error(t, err, "passing through required check")
	})

	t.Run("should update comment", func(t *testing.T) {
		var res comments.UpdateCommentHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v1/comments/%s", testData.TestUserOneCommentOneID),
			Body: &comments.UpdateCommentHandlerRequestBody{
				Name: updatedCommentName,
				Body: updatedCommentBody,
			},
			Response: &res,
		})
		require.NoError(t, err, "passing through required check")
		require.NotEmpty(t, res.Comment.ID, "id field was empty")
		require.NotEmpty(t, res.Comment.Name, "name field was empty")
		require.NotEmpty(t, res.Comment.Body, "body field was empty")
	})

	t.Run("comment updated in database", func(t *testing.T) {
		var comment models.Comment
		err := m.MySQL.
			Model(&models.Comment{}).
			Where(&models.Comment{Base: models.Base{ID: testData.TestUserOneCommentOneID}}).
			First(&comment).
			Error
		require.NoError(t, err, "failed to find comment in database")
		require.NotEmpty(t, comment.ID, "id field was empty")
		require.Equal(t, updatedCommentName, comment.Name, "unexpected name value")
		require.Equal(t, updatedCommentBody, comment.Body, "unexpected body value")
	})
}

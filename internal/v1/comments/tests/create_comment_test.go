package tests

import (
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v1/comments"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateCommentHandler(t *testing.T) {
	// init service
	m := app.Monolith{}
	m.Setup()

	// init test fixtures
	fixture := CommentsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	testClient := testclient.TestClient{}
	testClient.Setup(&testclient.Options{
		Router: m.Router,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, m.Config.HMACSecret),
	})

	unauthorizedClient := testclient.TestClient{}
	unauthorizedClient.Setup(&testclient.Options{
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

	t.Run("only authorized user can create a comment", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := unauthorizedClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				PostID: testData.TestPostOneID.String(),
				Name:   "test comment name",
				Body:   "test comment body",
			},
			Response: &res,
		})
		require.Error(t, err, "unauthorized user created comment")
	})

	t.Run("should error if passing invalid post uuid", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				PostID: "invalid uuid",
				Name:   "test comment name",
				Body:   "test comment body",
			},
			Response: &res,
		})
		require.Error(t, err, "parsed invalid uuid")
	})

	t.Run("post id field should be required", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				Name: "test comment name",
				Body: "test comment body",
			},
			Response: &res,
		})
		require.Error(t, err, "passed through required check")
	})

	t.Run("name field should be required", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				PostID: testData.TestPostOneID.String(),
				Body:   "test comment body",
			},
			Response: &res,
		})
		require.Error(t, err, "passed through required check")
	})

	t.Run("body field should be required", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				PostID: testData.TestPostOneID.String(),
				Name:   "test comment name",
			},
			Response: &res,
		})
		require.Error(t, err, "passed through required check")
	})

	var newCommentID uuid.UUID
	t.Run("comment should be created", func(t *testing.T) {
		var res comments.CreateCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/comments",
			Body: &comments.CreateCommentHandlerRequestBody{
				PostID: testData.TestPostOneID.String(),
				Name:   "test comment name",
				Body:   "test comment body",
			},
			Response: &res,
		})
		require.NoError(t, err, "failed to create comment")
		require.NotEmpty(t, res.Comment.ID, "id field was empty")
		require.NotEmpty(t, res.Comment.Name, "name field was empty")
		require.NotEmpty(t, res.Comment.Body, "body field was empty")
		newCommentID = res.Comment.ID
	})

	t.Run("comment should be saved in database", func(t *testing.T) {
		var comment models.Comment
		err := m.MySQL.
			Model(&models.Comment{}).
			Where(&models.Comment{Base: models.Base{ID: newCommentID}}).
			First(&comment).
			Error
		require.NoError(t, err, "failed to find comment in database")
		require.NotEmpty(t, comment.ID, "id field was empty")
		require.NotEmpty(t, comment.Name, "name field was empty")
		require.NotEmpty(t, comment.Body, "body field was empty")
	})
}

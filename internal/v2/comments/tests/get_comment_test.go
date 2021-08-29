package tests

import (
	"fmt"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/comments"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetCommentHandler(t *testing.T) {
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
		Router: m.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, m.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		var res comments.GetCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/comments/%s", "invalid uuid"),
			Response: &res,
		})
		require.Error(t, err, "parsed invalid uuid")
	})

	t.Run("should error if passing id of not existing comment", func(t *testing.T) {
		var res comments.GetCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/comments/%s", uuid.New().String()),
			Response: &res,
		})
		require.Error(t, err, "got not existing comment")
	})

	t.Run("should get requested comment", func(t *testing.T) {
		var res comments.GetCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/comments/%s", testData.TestUserOneCommentOneID),
			Response: &res,
		})
		require.NoError(t, err, "failed to get comment with provided id")
	})

	t.Run("fields should not be empty", func(t *testing.T) {
		var res comments.GetCommentHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/comments/%s", testData.TestUserOneCommentOneID),
			Response: &res,
		})
		require.NoError(t, err, "failed to get comment with provided id")
		require.NotEmpty(t, res.Comment.Name, "comment name was empty")
		require.NotEmpty(t, res.Comment.Body, "comment body was empty")
	})
}

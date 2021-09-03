package tests

import (
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/comments"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetCommentsHandler(t *testing.T) {
	// init service
	a := app.Application{}
	a.Setup()

	// init test fixtures
	fixture := CommentsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	testClient := testclient.TestClient{}
	testClient.Setup(&testclient.Options{
		Router: a.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, a.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	t.Run("should error if passing invalid user id", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
				UserID: "invalid uuid",
			},
			Response: &res,
		})
		require.Error(t, err, "parsed invalid uuid")
	})

	t.Run("should error if passing invalid post id", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
				PostID: "invalid uuid",
			},
			Response: &res,
		})
		require.Error(t, err, "parsed invalid uuid")
	})

	t.Run("should get all comments", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Equal(t, testData.TotalComments, int(res.Total), "invalid total length")
	})

	t.Run("should get all comments related to user", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
				UserID: testData.TestUserOneID.String(),
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Equal(t, testData.TotalCommentsByUserOne, int(res.Total), "invalid total length")
	})

	t.Run("should get all comments related to post", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
				PostID: testData.TestPostOneID.String(),
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Equal(t, testData.TotalCommentsInPostOne, int(res.Total), "invalid total length")
	})

	var prevCommentId uuid.UUID
	t.Run("limit should work", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  1,
				Offset: 0,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Len(t, *res.Comments, 1, "invalid response length")

		for _, c := range *res.Comments {
			prevCommentId = c.ID
		}
	})

	t.Run("offset should work", func(t *testing.T) {
		var res comments.GetCommentsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/comments",
			Query: &comments.GetCommentsHandlerRequestQuery{
				Limit:  1,
				Offset: 1,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")

		for _, c := range *res.Comments {
			require.NotEqual(t, prevCommentId, c.ID, "got same comment with different offset")
		}
	})

}

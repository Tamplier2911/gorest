package tests

import (
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/posts"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetPostsHandler(t *testing.T) {
	// init service
	m := app.Monolith{}
	m.Setup()

	// init test fixtures
	fixture := PostsTestFixtures()
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

	t.Run("should get all posts", func(t *testing.T) {
		var res posts.GetPostsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/posts",
			Query: &posts.GetPostsHandlerRequestQuery{
				Limit:  20,
				Offset: 0,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Equal(t, testData.TotalPosts, int(res.Total), "invalid total length")
	})

	var prevPostId uuid.UUID
	t.Run("limit should work", func(t *testing.T) {
		var res posts.GetPostsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/posts",
			Query: &posts.GetPostsHandlerRequestQuery{
				Limit:  1,
				Offset: 0,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")
		require.Len(t, *res.Posts, 1, "invalid response length")

		for _, p := range *res.Posts {
			prevPostId = p.ID
		}
	})

	t.Run("offset should work", func(t *testing.T) {
		var res posts.GetPostsHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "GET",
			URL:    "/api/v2/posts",
			Query: &posts.GetPostsHandlerRequestQuery{
				Limit:  1,
				Offset: 1,
			},
			Response: &res,
		})
		require.NoError(t, err, "unexpected response")

		for _, c := range *res.Posts {
			require.NotEqual(t, prevPostId, c.ID, "got same post with different offset")
		}
	})

}

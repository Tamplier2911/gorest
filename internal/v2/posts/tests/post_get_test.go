package tests

import (
	"fmt"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/posts"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetPostHandler(t *testing.T) {
	// init service
	a := app.Application{}
	a.Setup()

	// init test fixtures
	fixture := PostsTestFixtures()
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

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		var res posts.GetPostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/posts/%s", "invalid uuid"),
			Response: &res,
		})
		require.Error(t, err, "parsed invalid uuid")
	})

	t.Run("should error if passing id of not existing post", func(t *testing.T) {
		var res posts.GetPostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/posts/%s", uuid.New().String()),
			Response: &res,
		})
		require.Error(t, err, "got not existing post")
	})

	t.Run("should get requested post", func(t *testing.T) {
		var res posts.GetPostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Response: &res,
		})
		require.NoError(t, err, "failed to get posts with provided id")
	})

	t.Run("fields should not be empty", func(t *testing.T) {
		var res posts.GetPostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method:   "GET",
			URL:      fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Response: &res,
		})
		require.NoError(t, err, "failed to get post with provided id")
		require.NotEmpty(t, res.Post.Title, "post title was empty")
		require.NotEmpty(t, res.Post.Body, "post body was empty")
	})
}

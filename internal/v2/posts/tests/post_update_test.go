package tests

import (
	"fmt"
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v2/posts"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdatePostHandler(t *testing.T) {
	// init service
	a := app.Application{}
	a.Setup()

	// init test fixtures
	fixture := PostsTestFixtures()
	testData, err := fixture.Setup()
	require.NoError(t, err, "failed to setup test fixtures")

	// init test client
	authorClient := testclient.TestClient{}
	authorClient.Setup(&testclient.Options{
		Router: a.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: testData.TestUserOneID,
		}, a.Config.HMACSecret),
	})

	notAuthorClient := testclient.TestClient{}
	notAuthorClient.Setup(&testclient.Options{
		Router: a.Echo,
		Token: access.MustEncodeToken(&access.Token{
			UserID: uuid.New(),
		}, a.Config.HMACSecret),
	})

	defer func() {
		// cleanup test data
		err := fixture.Teardown()
		require.NoError(t, err, "failed to clean up test fixtures")
	}()

	updatedPostTitle := "updated test post name"
	updatedPostBody := "updated test post name"
	t.Run("only author can update post", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := notAuthorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Body: &posts.UpdatePostHandlerRequestBody{
				Title: updatedPostTitle,
				Body:  updatedPostBody,
			},
			Response: &res,
		})
		require.Error(t, err, "random user could update comment")
	})

	t.Run("should error if passing invalid post uuid", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", uuid.New()),
			Body: &posts.UpdatePostHandlerRequestBody{
				Title: updatedPostTitle,
				Body:  updatedPostBody,
			},
			Response: &res,
		})
		require.Error(t, err, "updated not existing comment in database")
	})

	t.Run("should error if passing invalid uuid", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", "invalid uuid"),
			Body: &posts.UpdatePostHandlerRequestBody{
				Title: updatedPostTitle,
				Body:  updatedPostBody,
			},
			Response: &res,
		})
		require.Error(t, err, "passing uuid parsing with invalid id")
	})

	t.Run("title field should be required", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Body: &posts.UpdatePostHandlerRequestBody{
				Body: updatedPostBody,
			},
			Response: &res,
		})
		require.Error(t, err, "passing through required check")
	})

	t.Run("body field should be required", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Body: &posts.UpdatePostHandlerRequestBody{
				Title: updatedPostTitle,
			},
			Response: &res,
		})
		require.Error(t, err, "passing through required check")
	})

	t.Run("should update post", func(t *testing.T) {
		var res posts.UpdatePostHandlerResponseBody
		err := authorClient.Request(&testclient.RequestOptions{
			Method: "PUT",
			URL:    fmt.Sprintf("/api/v2/posts/%s", testData.TestPostOneUserOneID),
			Body: &posts.UpdatePostHandlerRequestBody{
				Title: updatedPostTitle,
				Body:  updatedPostBody,
			},
			Response: &res,
		})
		require.NoError(t, err, "passing through required check")
		require.NotEmpty(t, res.Post.ID, "id field was empty")
		require.NotEmpty(t, res.Post.Title, "title field was empty")
		require.NotEmpty(t, res.Post.Body, "body field was empty")
	})

	t.Run("post updated in database", func(t *testing.T) {
		var post models.Post
		err := a.MySQL.
			Model(&models.Post{}).
			Where(&models.Post{Base: models.Base{ID: testData.TestPostOneUserOneID}}).
			First(&post).
			Error
		require.NoError(t, err, "failed to find post in database")
		require.NotEmpty(t, post.ID, "id field was empty")
		require.Equal(t, updatedPostTitle, post.Title, "unexpected title value")
		require.Equal(t, updatedPostBody, post.Body, "unexpected body value")
	})
}

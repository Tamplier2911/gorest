package tests

import (
	"testing"

	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/internal/v1/posts"
	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/testclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatePostHandler(t *testing.T) {
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

	t.Run("only authorized user can create a post", func(t *testing.T) {
		var res posts.CreatePostHandlerResponseBody
		err := unauthorizedClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/posts",
			Body: &posts.CreatePostHandlerRequestBody{
				Title: "test post title",
				Body:  "test post body",
			},
			Response: &res,
		})
		require.Error(t, err, "unauthorized user created post")
	})

	t.Run("title field should be required", func(t *testing.T) {
		var res posts.CreatePostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/posts",
			Body: &posts.CreatePostHandlerRequestBody{
				Body: "test post body",
			},
			Response: &res,
		})
		require.Error(t, err, "passed through required check")
	})

	t.Run("body field should be required", func(t *testing.T) {
		var res posts.CreatePostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/posts",
			Body: &posts.CreatePostHandlerRequestBody{
				Title: "test post title",
			},
			Response: &res,
		})
		require.Error(t, err, "passed through required check")
	})

	var newPostId uuid.UUID
	t.Run("post should be created", func(t *testing.T) {
		var res posts.CreatePostHandlerResponseBody
		err := testClient.Request(&testclient.RequestOptions{
			Method: "POST",
			URL:    "/api/v1/posts",
			Body: &posts.CreatePostHandlerRequestBody{
				Title: "test post title",
				Body:  "test post body",
			},
			Response: &res,
		})
		require.NoError(t, err, "failed to create post")
		require.NotEmpty(t, res.Post.ID, "id field was empty")
		require.NotEmpty(t, res.Post.Title, "title field was empty")
		require.NotEmpty(t, res.Post.Body, "body field was empty")
		newPostId = res.Post.ID
	})

	t.Run("post should be saved in database", func(t *testing.T) {
		var post models.Post
		err := m.MySQL.
			Model(&models.Post{}).
			Where(&models.Post{Base: models.Base{ID: newPostId}}).
			First(&post).
			Error
		require.NoError(t, err, "failed to find comment in database")
		require.NotEmpty(t, post.ID, "id field was empty")
		require.NotEmpty(t, post.Title, "title field was empty")
		require.NotEmpty(t, post.Body, "body field was empty")
	})
}

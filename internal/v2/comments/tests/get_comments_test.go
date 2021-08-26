package tests

import (
	"testing"

	"github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetCommentsHandler(t *testing.T) {

	// init service
	m := internal.Monolith{}

	// init test fixtures
	testUser := models.User{
		Username: "user_get_comments_test",
		Email:    "user_get_comments@test.com",
		UserRole: models.UserRoleUser,
	}
	err := m.MySQL.Create(&testUser).Error
	require.NoError(t, err, "failed to create test user")

	testPosts := []models.Post{
		{
			UserID: testUser.ID,
			Title:  "test post 1",
			Body:   "test post 1",
		},
		{
			UserID: uuid.New(),
			Title:  "test post 1",
			Body:   "test post 1",
		},
	}
	err = m.MySQL.Create(&testPosts).Error
	require.NoError(t, err, "failed to create test posts")

	testComments := []models.Comment{
		{
			PostID: testPosts[0].ID,
			UserID: testUser.ID,
			Name:   "comment 1",
			Body:   "comment 1",
		},
		{
			PostID: testPosts[0].ID,
			UserID: uuid.New(),
			Name:   "comment 2",
			Body:   "comment 2",
		},
		{
			PostID: testPosts[1].ID,
			UserID: testUser.ID,
			Name:   "comment 3",
			Body:   "comment 3",
		},
		{
			PostID: testPosts[1].ID,
			UserID: uuid.New(),
			Name:   "comment 4",
			Body:   "comment 4",
		},
		{
			PostID: testPosts[1].ID,
			UserID: uuid.New(),
			Name:   "comment 5",
			Body:   "comment 5",
		},
		{
			PostID: uuid.New(),
			UserID: uuid.New(),
			Name:   "comment 6",
			Body:   "comment 6",
		},
		{
			PostID: uuid.New(),
			UserID: uuid.New(),
			Name:   "comment 7",
			Body:   "comment 7",
		},
	}
	err = m.MySQL.Create(&testComments).Error
	require.NoError(t, err, "failed to create test posts")

	defer func() {
		// clean up test data
		err = m.MySQL.Unscoped().Delete(&testComments).Error
		require.NoError(t, err, "failed to delete test comments")

		err = m.MySQL.Unscoped().Delete(&testPosts).Error
		require.NoError(t, err, "failed to delete test posts")

		err := m.MySQL.Unscoped().Delete(&testUser).Error
		require.NoError(t, err, "failed to delete test users")
	}()

	t.Run("should error if passing invalid user id", func(t *testing.T) {

	})

	t.Run("should error if passing invalid post id", func(t *testing.T) {

	})

	t.Run("should get all comments", func(t *testing.T) {

	})

	t.Run("should get all comments related to user", func(t *testing.T) {

	})

	t.Run("should get all comments related to post", func(t *testing.T) {

	})

	t.Run("limit should work", func(t *testing.T) {

	})

	t.Run("offset should work", func(t *testing.T) {

	})

}

package tests

import (
	app "github.com/Tamplier2911/gorest/internal"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
)

// Fixtures represent test fixture.
type Fixture struct {
	Setup    func() (TestFixturesData, error)
	Teardown func() error
}

// TestFixtureData represent set of test fixture data.
type TestFixturesData struct {
	TestUserOneID           uuid.UUID
	TestUserTwoID           uuid.UUID
	TestPostOneID           uuid.UUID
	TestPostTwoID           uuid.UUID
	TestUserOneCommentOneID uuid.UUID
	TestUserTwoCommentOneID uuid.UUID

	TotalComments          int
	TotalCommentsInPostOne int
	TotalCommentsInPostTwo int
	TotalCommentsByUserOne int
	TotalCommentsByUserTwo int
}

// CommentTestFixtures return instance of fixture.
func CommentsTestFixtures() Fixture {
	// init service
	m := app.Monolith{}
	m.Setup()

	// test users
	var testUsers []models.User
	// test posts
	var testPosts []models.Post
	// test comments
	var testComments []models.Comment

	setup := func() (TestFixturesData, error) {
		// create test users
		testUsers = []models.User{
			{
				Username: "test_user_one_comments",
				Email:    "test_user_one_comments@test.com",
				UserRole: models.UserRoleUser,
			},
			{
				Username: "test_user_two_comments_test",
				Email:    "test_user_two_comments@test.com",
				UserRole: models.UserRoleUser,
			},
		}
		err := m.MySQL.Create(&testUsers).Error
		if err != nil {
			return TestFixturesData{}, err
		}

		// create test posts
		testPosts = []models.Post{
			{
				UserID: testUsers[0].ID,
				Title:  "test post 1",
				Body:   "test post 1",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 1",
				Body:   "test post 1",
			},
		}
		err = m.MySQL.Create(&testPosts).Error
		if err != nil {
			return TestFixturesData{}, err
		}

		// create test comments
		testComments = []models.Comment{
			{
				PostID: testPosts[0].ID,
				UserID: testUsers[0].ID,
				Name:   "comment 1",
				Body:   "comment 1",
			},
			{
				PostID: testPosts[0].ID,
				UserID: testUsers[1].ID,
				Name:   "comment 2",
				Body:   "comment 2",
			},
			{
				PostID: testPosts[0].ID,
				UserID: testUsers[1].ID,
				Name:   "comment 3",
				Body:   "comment 3",
			},
			{
				PostID: testPosts[1].ID,
				UserID: testUsers[0].ID,
				Name:   "comment 4",
				Body:   "comment 4",
			},
			{
				PostID: testPosts[1].ID,
				UserID: testUsers[0].ID,
				Name:   "comment 5",
				Body:   "comment 5",
			},
			{
				PostID: testPosts[1].ID,
				UserID: testUsers[1].ID,
				Name:   "comment 6",
				Body:   "comment 6",
			},
			{
				PostID: testPosts[1].ID,
				UserID: testUsers[1].ID,
				Name:   "comment 7",
				Body:   "comment 7",
			},
			{
				PostID: testPosts[1].ID,
				UserID: testUsers[1].ID,
				Name:   "comment 8",
				Body:   "comment 8",
			},
		}
		err = m.MySQL.Create(&testComments).Error
		if err != nil {
			return TestFixturesData{}, err
		}

		var totalByUserOne int
		var totalByUserTwo int
		var totalInPostOne int
		var totalInPostTwo int
		for _, c := range testComments {
			if c.UserID == testUsers[0].ID {
				totalByUserOne += 1
			}
			if c.UserID == testUsers[1].ID {
				totalByUserTwo += 1
			}
			if c.PostID == testPosts[0].ID {
				totalInPostOne += 1
			}
			if c.PostID == testPosts[1].ID {
				totalInPostTwo += 1
			}
		}

		return TestFixturesData{
			TestUserOneID:           testUsers[0].ID,
			TestUserTwoID:           testUsers[1].ID,
			TestPostOneID:           testPosts[0].ID,
			TestPostTwoID:           testPosts[1].ID,
			TestUserOneCommentOneID: testComments[0].ID,
			TestUserTwoCommentOneID: testComments[1].ID,

			TotalComments:          len(testComments),
			TotalCommentsByUserOne: totalByUserOne,
			TotalCommentsByUserTwo: totalByUserTwo,
			TotalCommentsInPostOne: totalInPostOne,
			TotalCommentsInPostTwo: totalInPostTwo,
		}, nil
	}

	teardown := func() error {
		// clean up test comments
		err := m.MySQL.Unscoped().Delete(&testComments).Error
		if err != nil {
			return err
		}

		// clean up test posts
		err = m.MySQL.Unscoped().Delete(&testPosts).Error
		if err != nil {
			return err
		}

		// clean up test users
		err = m.MySQL.Unscoped().Delete(&testUsers).Error
		if err != nil {
			return err
		}

		return nil
	}

	return Fixture{
		Setup:    setup,
		Teardown: teardown,
	}
}

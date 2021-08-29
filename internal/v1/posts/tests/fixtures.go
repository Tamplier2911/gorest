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
	TestUserOneID        uuid.UUID
	TestUserTwoID        uuid.UUID
	TestPostOneUserOneID uuid.UUID
	TestPostOneUserTwoID uuid.UUID

	TotalPosts int
}

// PostsTestFixtures return instance of fixture.
func PostsTestFixtures() Fixture {
	// init service
	m := app.Monolith{}
	m.Setup()

	// test users
	var testUsers []models.User
	// test posts
	var testPosts []models.Post

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
				UserID: testUsers[0].ID,
				Title:  "test post 2",
				Body:   "test post 2",
			},
			{
				UserID: testUsers[0].ID,
				Title:  "test post 3",
				Body:   "test post 3",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 4",
				Body:   "test post 4",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 5",
				Body:   "test post 5",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 6",
				Body:   "test post 6",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 7",
				Body:   "test post 7",
			},
			{
				UserID: testUsers[1].ID,
				Title:  "test post 8",
				Body:   "test post 8",
			},
		}
		err = m.MySQL.Create(&testPosts).Error
		if err != nil {
			return TestFixturesData{}, err
		}

		return TestFixturesData{
			TestUserOneID:        testUsers[0].ID,
			TestUserTwoID:        testUsers[1].ID,
			TestPostOneUserOneID: testPosts[0].ID,
			TestPostOneUserTwoID: testPosts[1].ID,

			TotalPosts: len(testPosts),
		}, nil
	}

	teardown := func() error {
		// clean up test posts
		err := m.MySQL.Unscoped().Delete(&testPosts).Error
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

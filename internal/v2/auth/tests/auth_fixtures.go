package tests

import (
	"time"

	"github.com/Tamplier2911/gorest/internal/v2/auth"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GooogleFixtures struct {
	Token      *oauth2.Token
	GoogleUser *auth.GoogleUserData
}

func GetGoogleFixtures() GooogleFixtures {
	return GooogleFixtures{
		Token: &oauth2.Token{
			AccessToken:  "super-secure-google-access-token",
			TokenType:    "Bearer",
			RefreshToken: "super-secure-google-refresh-token",
			Expiry:       time.Now().Add(time.Hour * 24 * 14),
		},
		GoogleUser: &auth.GoogleUserData{
			ID:            "Google#123",
			Email:         "google_auth@test.com",
			EmailVerified: true,
			Picture:       "https://picsum.photos/200/200",
		},
	}
}

type FacebookFixtures struct {
	Token        *oauth2.Token
	FacebookUser *auth.FacebookUserData
}

func GetFacebookFixtures() FacebookFixtures {
	return FacebookFixtures{
		Token: &oauth2.Token{
			AccessToken:  "super-secure-facebook-access-token",
			TokenType:    "Bearer",
			RefreshToken: "super-secure-facebook-refresh-token",
			Expiry:       time.Now().Add(time.Hour * 24 * 14),
		},
		FacebookUser: &auth.FacebookUserData{
			ID:    "Facebook#123",
			Email: "facebook_auth@test.com",
			Name:  "facebook auth",
		},
	}
}

type GithubFixtures struct {
	Token      *oauth2.Token
	GithubUser *github.User
}

func GetGithubFixtures() GithubFixtures {
	githubUid := int64(01234543210)
	githubEmail := "github_auth@test.com"
	githubName := "facebook auth"
	githubAvatar := "https://picsum.photos/300/300"

	return GithubFixtures{
		Token: &oauth2.Token{
			AccessToken:  "super-secure-github-access-token",
			TokenType:    "Bearer",
			RefreshToken: "super-secure-github-refresh-token",
			Expiry:       time.Now().Add(time.Hour * 24 * 14),
		},
		GithubUser: &github.User{
			ID:        &githubUid,
			Email:     &githubEmail,
			Name:      &githubName,
			AvatarURL: &githubAvatar,
		},
	}
}

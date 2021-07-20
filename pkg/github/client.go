package github

import (
	"context"
	"errors"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

func GetClient(user, pass, accessToken, sshFile string, appID, installationID int64, ctx context.Context) (*github.Client, error) {
	var httpClient *http.Client
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	if user != "" && pass != "" {
		tp := github.BasicAuthTransport{
			Username: strings.TrimSpace(user),
			Password: strings.TrimSpace(pass),
		}
		httpClient = tp.Client()
	}

	if sshFile != "" {
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, sshFile)
		if err != nil {
			return nil, err
		}
		httpClient = &http.Client{Transport: itr}
	}

	if httpClient == nil {
		return nil, errors.New("No GitHub credentials was supplied")
	}

	return github.NewClient(httpClient), nil
}

package github

import (
	"context"
	"github.com/google/go-github/v32/github"
)

type RepoService interface {
	AddAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.AdminEnforcement, *github.Response, error)
	GetAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.AdminEnforcement, *github.Response, error)
	RemoveAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
	UpdateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}
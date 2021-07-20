package tagChanger

import (
	"context"
	"errors"
	"github.com/OmerKahani/tagChanger/pkg/yamlChanger"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type RepositoryServiceStub struct {
	mock.Mock
}

func (r *RepositoryServiceStub) AddAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.AdminEnforcement, *github.Response, error) {
	panic("implement me")
}

func (r *RepositoryServiceStub) GetAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.AdminEnforcement, *github.Response, error) {
	args := r.Called(ctx, owner, repo, branch)

	return args.Get(0).(*github.AdminEnforcement), args.Get(1).(*github.Response), args.Error(2)

}

func (r *RepositoryServiceStub) RemoveAdminEnforcement(ctx context.Context, owner, repo, branch string) (*github.Response, error) {
	panic("implement me")
}

func (r *RepositoryServiceStub) GetContents(ctx context.Context, owner, repo, path string, opt *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error){
	args := r.Called(ctx, owner, repo, path, opt)

	return args.Get(0).(*github.RepositoryContent), args.Get(1).([]*github.RepositoryContent), args.Get(2).(*github.Response), args.Error(3)

}

func (r *RepositoryServiceStub) UpdateFile(ctx context.Context, owner, repo, path string, opt *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error){
	args := r.Called(ctx, owner, repo, path, opt)

	return args.Get(0).(*github.RepositoryContentResponse), args.Get(1).(*github.Response), args.Error(2)
}

func TestChangeFile(t *testing.T ) {
	t.Parallel()
	cases := []struct {
		name			string
		content      	string
		newContent   	string
		repo      	 	string
		branch     	 	string
		filePath     	string
		valuePath     	string
		newValue     	string
		expectedError	error
	}{

		{
			name: "happy flow handle yaml anchor",
			content : `
global: &global
 tag: 123

test: *global
`,
			newContent : `global: &global
    tag: "1234"
test: *global
`,
			repo: "owner/repository",
			branch: "branch",
			filePath: "filePath",
			valuePath: "global.tag",
			newValue: "1234",
		},
		{
			name: "happy flow",
			content : `
global:
 tag: "123"
`,
			newContent : `global:
    tag: "newValue"
`,
			repo: "owner/repository",
			branch: "branch",
			filePath: "filePath",
			valuePath: "global.tag",
			newValue: "newValue",

		},

		{
			name: "happy flow - other tags aren't changed",
			content : `
global:
 tag: 123
 other: other
other: other
`,
			newContent : `global:
    tag: "newValue"
    other: other
other: other
`,
			repo: "owner/repository",
			branch: "branch",
			filePath: "filePath",
			valuePath: "global.tag",
			newValue: "newValue",

		},

		{
			name: "happy flow number format as string",
			content : `
global:
 tag: 123
`,
			newContent : `global:
    tag: "1234"
`,
			repo: "owner/repository",
			branch: "branch",
			filePath: "filePath",
			valuePath: "global.tag",
			newValue: "1234",
		},

		{
			name: "repo format error",
			content : `
global:
 tag: 123
`,
			newContent : `global:
    tag: "newValue"
`,
			repo: "repo",
			branch: "branch",
			filePath: "filePath",
			valuePath: "global.tag",
			newValue: "newValue",
			expectedError: errors.New("--repo formant should be owner/repository"),
		},

		{
			name: "empty value path",
			content : `
global:
 tag: 123
`,
			newContent : `global:
    tag: "newValue"
`,
			repo: "owner/repository",
			branch: "branch",
			filePath: "filePath",
			valuePath: "",
			newValue: "newValue",
			expectedError: &yamlChanger.PathError{},
		},



	}


	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T){
		clientStub := RepositoryServiceStub{}
		ctx := context.Background()

		fileSHA := "file SHA"
		commitSHA := "commit SHA"
		clientStub.On("GetAdminEnforcement", ctx, "owner", "repository", "branch").
			Return( &github.AdminEnforcement{}, &github.Response{}, errors.New("error"))

		clientStub.On("GetContents", ctx, "owner", "repository", "filePath", &github.RepositoryContentGetOptions{
			Ref: c.branch,
		}).
			Return( &github.RepositoryContent{
				Content: &c.content,
				SHA: &fileSHA,
			}, []*github.RepositoryContent{}, &github.Response{}, nil)

		clientStub.On("UpdateFile", ctx, "owner", "repository", "filePath", &github.RepositoryContentFileOptions{
			Content: []byte(c.newContent),
			Message: github.String(""),
			SHA: &fileSHA,
			Branch: github.String(c.branch),
		}).
			Return( &github.RepositoryContentResponse{
				Commit: github.Commit{
					SHA: &commitSHA,
				},
			}, &github.Response{}, nil).Once()

		err := changeFile(context.Background(), &clientStub, c.repo, c.branch, c.filePath, c.valuePath, c.newValue)

		if c.expectedError != nil {
			assert.EqualError(t, err, c.expectedError.Error())
			return
		}

		if err != nil {
			t.Errorf("Got unxpeted error: %s", err.Error())
			return
		}

		clientStub.AssertExpectations(t)
		})
	}
}

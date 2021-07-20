package tagChanger

import (
	"context"
	"errors"
	"fmt"
	"github.com/OmerKahani/tagChanger/pkg/github"
	"github.com/OmerKahani/tagChanger/pkg/yamlChanger"
	goGithub "github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	filePath		string
	repo			string
	branch			string
	user			string
	pass			string
	accessToken		string
	valuePath		string
	newValue		string
	commitMessage	string
	sshFile			string
	appID			int64
	installationID	int64
)

func GetCommand() *cobra.Command{
	cmd := &cobra.Command{
		Use:	"tagChanger",
		Short:	"change yaml value in github",
		Long: 	"changer - retrieve a yaml file from github, change a specific value in it, and commit it back to github",
		RunE: 	func(_ *cobra.Command, _ []string) error{
			ctx := context.Background()

			client, err := github.GetClient(user, pass, accessToken, sshFile, appID, installationID, ctx)
			if err != nil {
				return err
			}

			err = changeFile(ctx, client.Repositories, repo, branch, filePath, valuePath, newValue)
			if err != nil {
				return err
			}

			return nil

		},
	}

	viper.AutomaticEnv()
	cmd.PersistentFlags().StringVar(&filePath,		"file-path",		"",						"the YAML file path")
	cmd.PersistentFlags().StringVar(&repo,			"repo",			"", 						"owner/repository")
	cmd.PersistentFlags().StringVar(&branch,		"branch",			"main",					"The github branch (default is main)")
	cmd.PersistentFlags().StringVar(&user,			"user",			"", 						"github user")
	cmd.PersistentFlags().StringVar(&pass,			"pass",			viper.GetString("PASS"),	"github password")
	cmd.PersistentFlags().StringVar(&accessToken,	"access-token",	viper.GetString("TOKEN"),	"github access token")
	cmd.PersistentFlags().StringVar(&sshFile,		"ssh-file",		"",						"github ssh key path (.pem file)")
	cmd.PersistentFlags().Int64Var(&appID,			"app-id",			0,						"the id of the application")
	cmd.PersistentFlags().Int64Var(&installationID,	"installation-id",0,						"the id of the application installation")
	cmd.PersistentFlags().StringVar(&valuePath,		"value-path",		"",						"the yaml path to the value")
	cmd.PersistentFlags().StringVar(&newValue,		"new-value",		"",						"the new value")
	cmd.PersistentFlags().StringVar(&commitMessage, "commit-msg",		"",						"the commit message that will be used")

	return cmd
}

func changeFile(ctx context.Context, client github.RepoService, repo, branch, filePath, valuePath, newValue string) error {
	repoSplits := strings.Split(repo, "/")
	if len(repoSplits) != 2 {
		return errors.New("--repo formant should be owner/repository")
	}

	file, _ , _, err := client.GetContents(ctx, repoSplits[0], repoSplits[1], filePath, &goGithub.RepositoryContentGetOptions{
		Ref: branch,
	})
	if err != nil {
		return err
	}

	fmt.Println("Got the file from github")

	decodedContent, err := file.GetContent()
	if err != nil {
		return err
	}

	body := yaml.Node{}
	err = yaml.Unmarshal([]byte(decodedContent), &body)
	if err != nil {
		return err
	}

	path, err := yamlChanger.GetPathSplits(valuePath)
	if err != nil {
		return err
	}

	err = yamlChanger.ChangeYaml(&body, newValue, path)
	if err != nil {
		return err
	}

	changedContent, err := yaml.Marshal(&body)
	if err != nil {
		return err
	}

	fmt.Println("changed yaml, now committing it to github")

	err = AdminForceDisable(ctx, client, repoSplits[0], repoSplits[1], branch,
		func() error {
			_,_, err := client.UpdateFile(ctx, repoSplits[0], repoSplits[1], filePath, &goGithub.RepositoryContentFileOptions{
				Branch:		&branch,
				Message: 	&commitMessage,
				Content: 	[]byte(changedContent),
				SHA:     	goGithub.String(file.GetSHA()),
			})

			return err
		})

	if err != nil {
		return err
	}

	fmt.Println("file committed")
	return nil
}
func AdminForceDisable(ctx context.Context, client github.RepoService, owner, repo, branch string, fn func() error) error {
	adminEnforc, _, adminErr := client.GetAdminEnforcement(ctx, owner, repo, branch)

	if adminErr != nil || (adminErr == nil && adminEnforc.Enabled == false) {
		return fn()
	}


	_, err := client.RemoveAdminEnforcement(ctx,owner, repo, branch)
	if err != nil {
		return err
	}

	FnErr := fn()

	_, _, err = client.AddAdminEnforcement(ctx,owner, repo, branch)
	if err != nil {
		return err
	}

	return FnErr
}
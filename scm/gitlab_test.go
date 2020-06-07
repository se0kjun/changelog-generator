package scm

import (
	"changelog-generator/config"
	"fmt"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestGitlabCommitAction(t *testing.T) {
	c := &config.Config{
		ScmConfig: config.ScmConfig{
			ScmType:           "gitlab",
			ScmAccessToken:    "<secret_access_token>",
			ScmApiBaseUrl:     "https://gitlab.com/api/v4/",
			ScmRepositoryInfo: "project_with_namespace",
			ScmPostAction: config.PostActionConfig{
				TargetBranch:  "master",
				AuthorEmail:   "hong921122@gmail.com",
				AuthorName:    "Seokjun Hong",
				CommitMessage: "Create README",
			},
		},
	}

	scmAction, _ := GetScmHandler(c)
	project, _ := scmAction.GetProject()
	gitlabProject := project.(*gitlab.Project)
	fmt.Println(gitlabProject)
	if err := scmAction.Commit(
		&gitlab.CommitAction{
			Action:   gitlab.FileUpdate,
			FilePath: "README.md",
			Content:  "README fffffffff",
		}, ""); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("commit succeed")
	}
}

func TestGetFiles(t *testing.T) {
	c := &config.Config{
		ScmConfig: config.ScmConfig{
			ScmType:           "gitlab",
			ScmAccessToken:    "<access_token>",
			ScmApiBaseUrl:     "https://gitlab.com/api/v4/",
			ScmRepositoryInfo: "projectId",
			ScmPostAction: config.PostActionConfig{
				TargetBranch:  "master",
				AuthorEmail:   "hong921122@gmail.com",
				AuthorName:    "Seokjun Hong",
				CommitMessage: "Create README",
			},
		},
	}

	scmAction, _ := GetScmHandler(c)
	files, err := scmAction.GetFiles("test/unreleased")
	if err == nil {
		for _, file := range files {
			fmt.Println(file.FilePath)
			fmt.Println(file.FileContent)
		}
	} else {
		fmt.Println(err)
	}
}

func TestGetFile(t *testing.T) {
	c := &config.Config{
		ScmConfig: config.ScmConfig{
			ScmType:           "gitlab",
			ScmAccessToken:    "<access_token>",
			ScmApiBaseUrl:     "https://gitlab.com/api/v4/",
			ScmRepositoryInfo: "projectId",
			ScmPostAction: config.PostActionConfig{
				TargetBranch:  "master",
				AuthorEmail:   "hong921122@gmail.com",
				AuthorName:    "Seokjun Hong",
				CommitMessage: "Create README",
			},
		},
	}

	scmAction, _ := GetScmHandler(c)
	file, err := scmAction.GetFile("README.md")
	if err == nil {
		fmt.Println(file.FilePath)
		fmt.Println(file.FileContent)
	} else {
		fmt.Println(err)
	}
}

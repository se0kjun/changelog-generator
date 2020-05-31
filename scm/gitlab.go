package scm

import (
	"changelog-generator/config"
	b64 "encoding/base64"
	"errors"

	"github.com/xanzy/go-gitlab"
)

type GitlabScm struct {
	conf              *config.Config
	repositoryUrl     string
	projectIdentifier interface{}
	targetBranch      string
	accessToken       string
	gitlabClient      *gitlab.Client
	gitlabProject     *gitlab.Project
}

func (g *GitlabScm) Init(c *config.Config) error {
	var err error
	g.projectIdentifier = c.GetRepositoryInfo()
	g.targetBranch = c.GetBranch()
	g.conf = c
	g.gitlabClient, err = gitlab.NewClient(c.GetAccessToken(), gitlab.WithBaseURL(c.GetApiBaseUrl()))
	g.gitlabProject, _, err = g.gitlabClient.Projects.GetProject(g.projectIdentifier, nil)
	return err
}

func (g *GitlabScm) Commit(item interface{}, message string) error {
	tmp := item.(*gitlab.CommitAction)
	authorEmail := g.conf.GetAuthorEmail()
	authorName := g.conf.GetAuthorName()
	branch := g.conf.GetBranch()
	commitMessage := g.conf.GetCommitMessage()
	opt := &gitlab.CreateCommitOptions{
		AuthorEmail:   &authorEmail,
		AuthorName:    &authorName,
		Branch:        &branch,
		CommitMessage: &commitMessage,
	}
	opt.Actions = append(opt.Actions, tmp)
	_, res, err := g.gitlabClient.Commits.CreateCommit(g.projectIdentifier, opt, nil)
	if res.StatusCode != 200 {
		return err
	}

	return nil
}

func (g *GitlabScm) Commits(item interface{}, message string) error {
	commitActions := item.([]*gitlab.CommitAction)
	authorEmail := g.conf.GetAuthorEmail()
	authorName := g.conf.GetAuthorName()
	branch := g.conf.GetBranch()
	commitMessage := g.conf.GetCommitMessage()
	opt := &gitlab.CreateCommitOptions{
		AuthorEmail:   &authorEmail,
		AuthorName:    &authorName,
		Branch:        &branch,
		CommitMessage: &commitMessage,
		Actions:       commitActions,
	}
	_, res, err := g.gitlabClient.Commits.CreateCommit(g.projectIdentifier, opt, nil)
	if res.StatusCode != 200 {
		return err
	}

	return nil
}

func (g *GitlabScm) TagList() ([]string, error) {
	if g.gitlabProject != nil {
		return g.gitlabProject.TagList, nil
	} else {
		return []string{}, errors.New("Not initialized")
	}
}

func (g *GitlabScm) GetProject() (interface{}, error) {
	return g.gitlabProject, nil
}

func (g *GitlabScm) GetFiles(path string) ([]ScmFile, error) {
	scmFiles := make([]ScmFile, 1)
	opt := &gitlab.ListTreeOptions{
		Path: &path,
		Ref:  &g.targetBranch,
	}

	nodes, _, err := g.gitlabClient.Repositories.ListTree(g.projectIdentifier, opt, nil)
	if err != nil {
		return nil, err
	}
	for _, file := range nodes {
		fileOpt := &gitlab.GetFileOptions{
			Ref: &g.targetBranch,
		}

		fileObj, _, err := g.gitlabClient.RepositoryFiles.GetFile(g.projectIdentifier, file.Path, fileOpt, nil)
		if err == nil {
			content, _ := b64.StdEncoding.DecodeString(fileObj.Content)

			scmFiles = append(scmFiles, ScmFile{
				FilePath:    fileObj.FilePath,
				FileName:    fileObj.FileName,
				FileContent: string(content),
			})
		} else {
			return nil, err
		}
	}

	return scmFiles, nil
}

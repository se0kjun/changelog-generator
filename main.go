package main

import (
	"changelog-generator/config"
	"changelog-generator/handler/change"
	"changelog-generator/markdown"
	"changelog-generator/scm"
	"io/ioutil"
	"os"

	gitlab "github.com/xanzy/go-gitlab"
)

func makeOutputAction(c *config.Config, m *markdown.MarkdownGenerator) error {
	return nil
}

func doLocalFilePostAction(c *config.Config, clh change.ChangeLogBuilder, result string) error {
	if err := ioutil.WriteFile(c.GetOutputFilePath(), []byte(result), 0644); err != nil {
		return err
	}

	if c.ScmConfig.ScmPostAction.RemoveChangeLogFiles {
		removedLogFiles := make([]string, 10)
		for _, logs := range clh.GetChangeLogInfo() {
			for _, item := range logs {
				if err := os.Remove(item.GetFilePath()); err != nil {
					panic(err)
				} else {
					removedLogFiles = append(removedLogFiles, item.GetFilePath())
				}
			}
		}
	}

	return nil
}

func doScmPostAction(c *config.Config, clh change.ChangeLogBuilder, result string) error {
	if scmHandler, err := scm.GetScmHandler(c); err != nil {
		return err
	} else {
		commits := make([]gitlab.CommitAction, 10)
		if c.ScmConfig.ScmPostAction.PushChangeLog {
			commits = append(commits, gitlab.CommitAction{
				FilePath: c.GetChangeLogPath(),
				Content:  result,
				Action:   gitlab.FileUpdate,
			})
		}
		if c.ScmConfig.ScmPostAction.PushRemovedFiles {
			for _, val := range clh.GetChangeLogInfo() {
				for _, item := range val {
					commits = append(commits, gitlab.CommitAction{
						FilePath: item.GetFilePath(),
						Action:   gitlab.FileDelete,
					})
				}
			}
		}
		if commitErr := scmHandler.Commits(commits, ""); commitErr != nil {
			return commitErr
		}
	}

	return nil
}

func main() {
	if c, err := config.LoadChangeLogConfig(""); err != nil {
		panic(err)
	} else {
		if handler, err := change.NewChangeLogHandler(c); err != nil {
			panic(err)
		} else {
			// t.Log(handler.ChangeLogs)
			if markdownGen, err := markdown.NewMarkdownGenerator(c, handler); err != nil {
				panic(err)
			} else {
				if str, err := markdownGen.MakeResult(); err != nil {
					panic(err)
				} else {
					switch c.GetProjectAccessType() {
					case config.PROJECT_ACCESS_SCM:
						postActionErr := doScmPostAction(c, handler, str)
						if postActionErr != nil {
							panic(postActionErr)
						}
						break
					case config.PROJECT_ACCESS_LOCALFILE:
						postActionErr := doLocalFilePostAction(c, handler, str)
						if postActionErr != nil {
							panic(postActionErr)
						}
						break
					}
				}
			}
		}
	}
}

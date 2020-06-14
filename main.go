package main

import (
	"changelog-generator/config"
	"changelog-generator/handler/change"
	"changelog-generator/markdown"
	"changelog-generator/scm"
	"flag"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
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
		commits := make([]*gitlab.CommitAction, 0)
		if c.ScmConfig.ScmPostAction.PushChangeLog {
			action := gitlab.FileUpdate
			_, err := scmHandler.GetFile(c.GetOutputFilePath())
			if err != nil {
				log.Errorf("getting file error: %s, %s", c.GetOutputFilePath(), err.Error())
				action = gitlab.FileCreate
			}
			commits = append(commits, &gitlab.CommitAction{
				FilePath: c.GetOutputFilePath(),
				Content:  result,
				Action:   action,
			})
		}
		if c.ScmConfig.ScmPostAction.PushRemovedFiles {
			for _, val := range clh.GetChangeLogInfo() {
				for _, item := range val {
					commits = append(commits, &gitlab.CommitAction{
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
	// load json configuration file or load flags
	var configFile *string
	configFile = flag.String("config", "", "json configuration file")
	flag.Parse()
	if c, err := config.LoadChangeLogConfig(*configFile); err != nil {
		log.Fatalf("fail to load configuration: %s", err)
	} else {
		// initialize changelog handler
		if handler, err := change.NewChangeLogHandler(c); err != nil {
			log.Fatalf("fail to initialize change log handler: %s", err)
		} else {
			// initialize markdown generator
			if markdownGen, err := markdown.NewMarkdownGenerator(c, handler); err != nil {
				log.Fatalf("fail to initialize markdown generator: %s", err)
			} else {
				// make changelog output to formatted markdown
				if str, err := markdownGen.MakeResult(); err != nil {
					log.Fatalf("fail to make markdown result: %s", err)
				} else {
					// check project access type
					switch c.GetProjectAccessType() {
					case config.PROJECT_ACCESS_GITLAB:
						// if project access type is SCM, markdown file should be executed SCM post action
						// if SCM post action has specified, it would be pushed markdown file to remote repository
						postActionErr := doScmPostAction(c, handler, str)
						if postActionErr != nil {
							log.Fatalf("fail to scm post action: %s", postActionErr)
						}
						break
					case config.PROJECT_ACCESS_LOCALFILE:
						// if project access type is LOCAL_FILE, it should be executed local file post action
						// save as local file
						postActionErr := doLocalFilePostAction(c, handler, str)
						if postActionErr != nil {
							log.Fatalf("fail to local post action: %s", postActionErr)
						}
						break
					}
					log.Infof("generation following content: %s", str)
				}
			}
		}
		log.Infof("summary: %s changelog file successfully updated", c.GetOutputFilePath())
	}
}

package main

import (
	"changelog-generator/config"
	"changelog-generator/handler/change"
	"changelog-generator/markdown"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestRun(t *testing.T) {
	// load json configuration file or load flags
	if c, err := config.LoadChangeLogConfig("./test/scm_config.json"); err != nil {
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

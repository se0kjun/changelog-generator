package main

import (
	"changelog-generator/config"
	"changelog-generator/handler/change"
	"changelog-generator/markdown"
	"testing"
)

func TestRun(t *testing.T) {
	// load json configuration file or load flags
	if c, err := config.LoadChangeLogConfig("./test/scm_config.json"); err != nil {
		panic(err)
	} else {
		// initialize changelog handler
		if handler, err := change.NewChangeLogHandler(c); err != nil {
			panic(err)
		} else {
			// initialize markdown generator
			if markdownGen, err := markdown.NewMarkdownGenerator(c, handler); err != nil {
				panic(err)
			} else {
				// make changelog output to formatted markdown
				if str, err := markdownGen.MakeResult(); err != nil {
					panic(err)
				} else {
					// check project access type
					switch c.GetProjectAccessType() {
					case config.PROJECT_ACCESS_GITLAB:
						// if project access type is SCM, markdown file should be executed SCM post action
						// if SCM post action has specified, it would be pushed markdown file to remote repository
						postActionErr := doScmPostAction(c, handler, str)
						if postActionErr != nil {
							panic(postActionErr)
						}
						break
					case config.PROJECT_ACCESS_LOCALFILE:
						// if project access type is LOCAL_FILE, it should be executed local file post action
						// save as local file
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

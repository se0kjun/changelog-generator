package change

import (
	"changelog-generator/config"
	changelog_err "changelog-generator/errors"

	log "github.com/sirupsen/logrus"
)

type ChangeLogBuilder interface {
	init(*config.Config) error
	collectLogs() error
	GetChangeLogInfo() map[string][]DefaultChangeLog
}

var ChangeLogMakerMap = map[string]ChangeLogBuilder{
	config.PROJECT_ACCESS_LOCALFILE: &LocalFileChangeLogHandler{},
	config.PROJECT_ACCESS_GITLAB:    &GitlabChangeLogHandler{},
}

/*
{
	"type": "added" | "updated",
	"description": "description",
	"author": "author",
	"pr": "1234",
	"issue": "1234",
	"version": "1.2.3",
	"userDefined": {
		"blahblah": {
			...
		}
	}
}
*/
type DefaultChangeLog struct {
	fileName  string
	extension string
	path      string
	/* required */
	ChangeType        string `json:"type"`
	ChangeDescription string `json:"description"`
	Author            string `json:"author"`
	/* additional */
	ChangeTitle string `json:"title"`
	PrNum       string `json:"pr"`
	IssueNumber string `json:"issue"`
	Version     string `json:"version"`
	/* optional */
	UserDefinedData map[string]interface{} `json:"userDefined"`
}

func NewChangeLogHandler(c *config.Config) (ChangeLogBuilder, error) {
	log.Infof("project access type is %s", c.GetProjectAccessType())
	if _, ok := ChangeLogMakerMap[c.GetProjectAccessType()]; !ok {
		return nil, changelog_err.NOT_FOUND_ACCESS_TYPE
	}

	if err := ChangeLogMakerMap[c.GetProjectAccessType()].init(c); err != nil {
		log.Errorf("changelog handler initialization failed: %s", err.Error())
		return nil, err
	}

	logHandler := ChangeLogMakerMap[c.GetProjectAccessType()]
	if err := logHandler.collectLogs(); err != nil {
		log.Errorf("collect changelog failed: %s", err.Error())
		return nil, err
	}

	return logHandler, nil
}

func (c *DefaultChangeLog) GetFileName() string {
	return c.fileName
}

func (c *DefaultChangeLog) GetFilePath() string {
	return c.path
}

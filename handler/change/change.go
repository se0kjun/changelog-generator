package change

import (
	"changelog-generator/config"
)

type ChangeLogBuilder interface {
	init(*config.Config) error
	collectLogs() error
	GetChangeLogInfo() map[string][]DefaultChangeLog
}

var ChangeLogMakerMap = map[string]ChangeLogBuilder{
	"localfile": &LocalFileChangeLogHandler{},
	"gitlab":    &GitlabChangeLogHandler{},
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
	if err := ChangeLogMakerMap[c.GetProjectAccessType()].init(c); err != nil {
		return nil, err
	}

	logHandler := ChangeLogMakerMap[c.GetProjectAccessType()]
	if err := logHandler.collectLogs(); err != nil {
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

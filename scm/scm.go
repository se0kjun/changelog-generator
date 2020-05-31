package scm

import "changelog-generator/config"

var ScmHandlerMap = map[string]ScmAction{
	"gitlab": &GitlabScm{},
}

const (
	DEFAULT_COMMIT_MESSAGE = ``
)

type ScmAction interface {
	Init(c *config.Config) error
	Commit(interface{}, string) error
	Commits(interface{}, string) error
	TagList() ([]string, error)
	GetProject() (interface{}, error)
	GetFiles(string) ([]ScmFile, error)
}

type ScmFile struct {
	FilePath    string
	FileName    string
	FileContent string
}

func initializeScmHandlerMap(c *config.Config) error {
	for _, val := range ScmHandlerMap {
		if err := val.Init(c); err != nil {
			return err
		}
	}

	return nil
}

func GetScmHandler(c *config.Config) (ScmAction, error) {
	if err := initializeScmHandlerMap(c); err != nil {
		return nil, err
	} else {
		return ScmHandlerMap[c.GetScmType()], nil
	}
}

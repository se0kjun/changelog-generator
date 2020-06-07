package project

import (
	"changelog-generator/config"

	log "github.com/sirupsen/logrus"
)

type ChangeLogHeader interface {
	Init(*config.Config) error
	GetProjectName() string
	GetProjectDescription() string
}

var ChangeLogHeaderHandler = map[string]ChangeLogHeader{
	"default": &ConfigChangeLogHeader{},
	// "scm": &ScmChangeLogHeader{},
}

type ConfigChangeLogHeader struct {
	projectName        string
	projectDescription string
}

func (p *ConfigChangeLogHeader) Init(c *config.Config) error {
	p.projectDescription = c.GetProjectDescription()
	p.projectName = c.GetProjectName()
	log.Infof("project name: %s", p.projectName)
	log.Infof("project description: %s", p.projectDescription)
	return nil
}

func (p *ConfigChangeLogHeader) GetProjectName() string {
	return p.projectName
}

func (p *ConfigChangeLogHeader) GetProjectDescription() string {
	return p.projectDescription
}

// type ScmChangeLogHeader struct {
// }

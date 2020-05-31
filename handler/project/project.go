package project

import "changelog-generator/config"

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

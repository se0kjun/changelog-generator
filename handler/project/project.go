package project

import (
	"changelog-generator/config"
	changelog_err "changelog-generator/errors"
	"changelog-generator/scm"

	"github.com/xanzy/go-gitlab"
)

type ChangeLogHeader interface {
	Init(*config.Config) error
	GetProjectName() string
	GetProjectDescription() string
}

var ChangeLogHeaderHandler = map[string]ChangeLogHeader{
	config.PROJECT_ACCESS_LOCALFILE: &ConfigChangeLogHeader{},
	config.PROJECT_ACCESS_GITLAB:    &ScmChangeLogHeader{},
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

type ScmChangeLogHeader struct {
	projectName        string
	projectDescription string
	gitlabScmObject    scm.ScmAction
}

func (s *ScmChangeLogHeader) Init(c *config.Config) error {
	var err error
	if s.gitlabScmObject, err = scm.GetScmHandler(c); err != nil {
		return err
	}

	if tmp, err := s.gitlabScmObject.GetProject(); err != nil {
		return err
	} else {
		switch p := tmp.(type) {
		case *gitlab.Project:
			s.projectName = p.Name
			s.projectDescription = p.Description
			break
		default:
			return changelog_err.UNKNOWN_PROJECT_TYPE
		}
	}

	return nil
}

func (s *ScmChangeLogHeader) GetProjectName() string {
	return s.projectName
}

func (s *ScmChangeLogHeader) GetProjectDescription() string {
	return s.projectDescription
}

package config

import (
	"flag"
)

/*
{
	"generate": {
		"writePolicy": "prepend" | "append",
		"noExist": "create" | "error",
		"originFile": "filename",
		"outputPath": "path"
	},
	"changelog": {
		"changelogPath": "path",
		"format": "file",
		"version": {
			"type": "filename" | "scm-tag",
			"rule": "^\\d+\\.\\d+\\.\\d"
		}
	},
	"header": {
		"format": "file",
		"project": "project",
		"description": "description"
	},
	"versioning" : {
		"version": "file",
		"replaceVersionTag": "<$version%>"
	},
	"scm": {
		"type": "gitlab",
		"repository": "url",
		"accessToken": "token",
		"projectId": "id",
		"author": {
			"email": "email address",
			"name": "name"
		},
		"postAction": {
			"targetBranch": "branch",
			"action": "commit",
			"policy": {
				"removeLogs": true
			},
			"commitMessage": {
			}
		}
	}
}
*/
type Config struct {
	// local or scm
	ProjectAccessType string
	ProjectBasePath   string
	OutputConfig      OutputConfig
	ChangeLogConfig   ChangeLogGenerateConfig
	HeaderConfig      HeaderGenerateConfig
	ScmConfig         ScmConfig
}

type (
	OutputConfig struct {
		// stdout, file, scm
		OutputType     string
		OutputFilePath string
		WritePolicy    string
	}

	ChangeLogGenerateConfig struct {
		ChangeLogOnly               bool
		ChangeLogPath               string
		MarkdownChangelogFormatFile string
		VersionParseRule            string
		VersionAcquisitionPolicy    string
	}

	HeaderGenerateConfig struct {
		MarkdownHeaderFormatFile string
		ProjectName              string
		ProjectDescription       string
	}

	ScmConfig struct {
		ScmType           string
		ScmRepositoryInfo interface{}
		ScmApiBaseUrl     string
		ScmAccessToken    string
		ScmPostAction     PostActionConfig
	}

	PostActionConfig struct {
		RemoveChangeLogFiles bool
		PushChangeLog        bool
		PushRemovedFiles     bool
		AuthorEmail          string
		AuthorName           string
		TargetBranch         string
		CommitMessage        string
		// AutomaticVersioning  string
	}
)

/* project access type */
const (
	PROJECT_ACCESS_SCM       = "scm"
	PROJECT_ACCESS_LOCALFILE = "localfile"
)

/* output file policy */
const (
	PREPEND = "prepend"
	APPEND  = "append"
	CREATE  = "create"
)

/* VersionAcquisitionPolicy */
const (
	VERSION_GET_FILENAME = "filename"
	VERSION_GET_SCM_TAG  = "scm-tag"
)

func LoadChangeLogConfig(file string) (*Config, error) {
	conf := &Config{}
	if file != "" {
		conf.loadConfigByJson(file)
		return conf, nil
	}

	flag.StringVar(&conf.ChangeLogConfig.ChangeLogPath, "path", "", "change log information")
	flag.StringVar(&conf.ChangeLogConfig.MarkdownChangelogFormatFile, "changeformat", "", "markdown changelog format file")
	flag.StringVar(&conf.ChangeLogConfig.VersionParseRule, "versionrule", "^\\d+\\.\\d+\\.\\d", "version parsing rule")
	flag.StringVar(&conf.ChangeLogConfig.VersionAcquisitionPolicy, "versionacquisition", "filename", "version parsing rule")
	flag.BoolVar(&conf.ChangeLogConfig.ChangeLogOnly, "changelog-only", false, "generate changelog only")

	flag.StringVar(&conf.HeaderConfig.MarkdownHeaderFormatFile, "headerformat", "", "markdown header format file")
	flag.StringVar(&conf.HeaderConfig.ProjectDescription, "project-description", "", "project description")
	flag.StringVar(&conf.HeaderConfig.ProjectName, "project-name", "", "project name")

	flag.StringVar(&conf.OutputConfig.OutputFilePath, "o", "./CHANGELOG.md", "generate markdown file on the specified path")

	flag.Parse()

	return conf, nil
}

func (c *Config) loadConfigByJson(file string) {
}

func (c *Config) GetProjectAccessType() string {
	return c.ProjectAccessType
}

/* output config */
func (c *Config) GetOutputFilePath() string {
	return c.OutputConfig.OutputFilePath
}

func (c *Config) GetWritePolicy() string {
	return c.OutputConfig.WritePolicy
}

/* change log config */
func (c *Config) GetChangeLogPath() string {
	return c.ChangeLogConfig.ChangeLogPath
}

func (c *Config) GetChangeLogFormat() string {
	return c.ChangeLogConfig.MarkdownChangelogFormatFile
}

func (c *Config) GetVersionParsingRule() string {
	return c.ChangeLogConfig.VersionParseRule
}

func (c *Config) GetVersionAcquisitionPolicy() string {
	return c.ChangeLogConfig.VersionAcquisitionPolicy
}

func (c *Config) IsChangeLogOnly() bool {
	return c.ChangeLogConfig.ChangeLogOnly
}

/* header config */
func (c *Config) GetHeaderFormat() string {
	return c.HeaderConfig.MarkdownHeaderFormatFile
}

func (c *Config) GetProjectName() string {
	return c.HeaderConfig.ProjectName
}

func (c *Config) GetProjectDescription() string {
	return c.HeaderConfig.ProjectDescription
}

/* SCM config */
func (c *Config) GetScmType() string {
	return c.ScmConfig.ScmType
}

func (c *Config) GetAccessToken() string {
	return c.ScmConfig.ScmAccessToken
}

func (c *Config) GetApiBaseUrl() string {
	return c.ScmConfig.ScmApiBaseUrl
}

func (c *Config) GetRepositoryInfo() interface{} {
	return c.ScmConfig.ScmRepositoryInfo
}

func (c *Config) GetAuthorEmail() string {
	return c.ScmConfig.ScmPostAction.AuthorEmail
}

func (c *Config) GetAuthorName() string {
	return c.ScmConfig.ScmPostAction.AuthorName
}

func (c *Config) GetBranch() string {
	return c.ScmConfig.ScmPostAction.TargetBranch
}

func (c *Config) GetCommitMessage() string {
	return c.ScmConfig.ScmPostAction.CommitMessage
}

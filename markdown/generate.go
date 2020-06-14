package markdown

import (
	"bytes"
	"changelog-generator/config"
	changelog_err "changelog-generator/errors"
	"changelog-generator/handler/change"
	"changelog-generator/handler/project"
	"changelog-generator/scm"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

/* default */
const (
	DEFAULT_MARKDOWN_CHANGELOG_TEMPLATE = `
{{ range $version, $section := .ChangeLogSection }}### ChangeLog ({{ $version }})
{{ range $type, $changelog := $section.ChangeLog }}
- {{ $type }}
	{{ range $changelog }}
	- {{ .Author }}: {{ .ChangeDescription }} (!{{ .ChangePr }}){{ end }}
{{ end }}
{{ end }}`
	DEFAULT_MARKDOWN_HEADER_TEMPLATE = `# {{ .Project }}

{{ .ProjectDescription }}

---
	`
)

type MarkdownGenerator struct {
	conf *config.Config
	// change log markdown data
	markdownPath      string
	markdownTemplate  string
	markdownParseTree *template.Template
	// change log header markdown data
	headerMarkdownPath      string
	headerMarkdownTemplate  string
	headerMarkdownParseTree *template.Template
	// data
	templateData           *MarkdownMapperInterface
	changeLogHandler       change.ChangeLogBuilder
	changeLogHeaderHandler project.ChangeLogHeader
}

type (
	/* markdown template toplevel interface */
	MarkdownMapperInterface struct {
		Project            string
		ProjectDescription string
		/* key: version, value: change logs */
		ChangeLogSection map[string]MarkdownChangelogSection
		DateTime         string
	}

	/* changelog section data */
	MarkdownChangelogSection struct {
		/* key: change type, value: change log */
		ChangeLog map[string][]MarkdownChangelog
		Version   string
		DateTime  string
		Date      string
	}

	MarkdownChangelog struct {
		ChangeType        string
		ChangeDescription string
		ChangePr          string
		ChangeIssue       string
		Author            string
		DateTime          string
		Date              string
		UserDefined       map[string]interface{}
	}
)

func NewMarkdownGenerator(c *config.Config, lh change.ChangeLogBuilder) (*MarkdownGenerator, error) {
	gen := new(MarkdownGenerator)
	gen.conf = c
	gen.markdownPath = c.GetChangeLogFormat()
	gen.headerMarkdownPath = c.GetHeaderFormat()
	if gen.markdownPath == "" {
		log.Infof("Use default changelog markdown template")
		gen.markdownTemplate = DEFAULT_MARKDOWN_CHANGELOG_TEMPLATE
	}
	if gen.headerMarkdownPath == "" {
		log.Infof("Use default header markdown template")
		gen.headerMarkdownTemplate = DEFAULT_MARKDOWN_HEADER_TEMPLATE
	}
	gen.SetChangeLogHandler(lh)

	if ch, ok := project.ChangeLogHeaderHandler[c.GetProjectAccessType()]; ok {
		gen.changeLogHeaderHandler = ch
	} else {
		return nil, changelog_err.NOT_FOUND_ACCESS_TYPE
	}

	if err := gen.changeLogHeaderHandler.Init(c); err != nil {
		log.Errorf("header handler initializing failed: %s", err.Error())
		return nil, err
	}
	if err := gen.init(); err != nil {
		log.Errorf("markdown generator initializing failed: %s", err.Error())
		return nil, err
	}

	return gen, nil
}

func (m *MarkdownGenerator) init() error {
	var err error
	if m.markdownPath != "" {
		if data, err := ioutil.ReadFile(m.markdownPath); err != nil {
			return err
		} else {
			m.markdownTemplate = string(data)
		}
	}

	if m.markdownParseTree, err = template.New("changelog").Parse(m.markdownTemplate); err != nil {
		return err
	}

	if m.headerMarkdownPath != "" {
		if data, err := ioutil.ReadFile(m.headerMarkdownPath); err != nil {
			return err
		} else {
			m.headerMarkdownTemplate = string(data)
		}
	}

	if m.headerMarkdownParseTree, err = template.New("changelog-header").Parse(m.headerMarkdownTemplate); err != nil {
		return err
	}

	m.buildDataMapper()

	return nil
}

func (m *MarkdownGenerator) buildDataMapper() {
	m.templateData = &MarkdownMapperInterface{
		ChangeLogSection: make(map[string]MarkdownChangelogSection),
	}
	m.templateData.Project = m.changeLogHeaderHandler.GetProjectName()
	m.templateData.ProjectDescription = m.changeLogHeaderHandler.GetProjectDescription()
	for version, logs := range m.changeLogHandler.GetChangeLogInfo() {
		logSection := MarkdownChangelogSection{
			Version:   version,
			ChangeLog: make(map[string][]MarkdownChangelog),
		}
		for _, logItem := range logs {
			item := MarkdownChangelog{
				ChangeType:        logItem.ChangeType,
				ChangeDescription: logItem.ChangeDescription,
				ChangeIssue:       logItem.IssueNumber,
				ChangePr:          logItem.PrNum,
				UserDefined:       logItem.UserDefinedData,
				Author:            logItem.Author,
			}
			if _, ok := logSection.ChangeLog[logItem.ChangeType]; ok {
				logSection.ChangeLog[logItem.ChangeType] = append(logSection.ChangeLog[logItem.ChangeType], item)
			} else {
				logSection.ChangeLog[logItem.ChangeType] = []MarkdownChangelog{item}
			}
		}

		m.templateData.ChangeLogSection[version] = logSection
	}
}

func (m *MarkdownGenerator) generateChangeLog() (string, error) {
	buf := new(bytes.Buffer)
	if err := m.markdownParseTree.Execute(buf, m.templateData); err != nil {
		return "", err
	} else {
		return buf.String(), nil
	}
}

func (m *MarkdownGenerator) generateHeader() (string, error) {
	buf := new(bytes.Buffer)
	if err := m.headerMarkdownParseTree.Execute(buf, m.templateData); err != nil {
		return "", err
	} else {
		return buf.String(), nil
	}
}

func (m *MarkdownGenerator) generate() (string, error) {
	var changeLogStr string
	var err error
	if m.conf.IsChangeLogOnly() {
		changeLogStr, err = m.generateChangeLog()
		if err != nil {
			return "", err
		}
	}

	headerStr, _ := m.generateHeader()
	return headerStr + changeLogStr, nil
}

func (m *MarkdownGenerator) MakeResult() (string, error) {
	var fileContent string
	var buffer bytes.Buffer

	log.Infof("project access type is %s", m.conf.GetProjectAccessType())
	switch m.conf.GetProjectAccessType() {
	case config.PROJECT_ACCESS_LOCALFILE:
		if _, err := os.Stat(m.conf.GetOutputFilePath()); os.IsNotExist(err) {
			return m.generate()
		}

		data, err := ioutil.ReadFile(m.conf.GetOutputFilePath())
		if err == nil {
			fileContent = string(data)
		} else {
			return "", err
		}
		break
	case config.PROJECT_ACCESS_GITLAB:
		scmAction, err := scm.GetScmHandler(m.conf)
		if err != nil {
			return "", err
		}
		file, fileErr := scmAction.GetFile(m.conf.GetOutputFilePath())
		if fileErr != nil {
			// file doesn't exist
			return m.generate()
		}
		fileContent = file.FileContent
		break
	default:
		return "", changelog_err.NOT_FOUND_ACCESS_TYPE
	}

	idx := strings.Index(fileContent, "---")
	if idx < 0 {
		return "", changelog_err.NOT_FOUND_DELIMETER
	}
	idx += len("---")

	log.Infof("output file policy: %s", m.conf.GetWritePolicy())
	markdownResult, _ := m.generateChangeLog()
	switch m.conf.GetWritePolicy() {
	case config.PREPEND:
		buffer.WriteString(fileContent[:idx])
		buffer.WriteString(markdownResult)
		buffer.WriteString(fileContent[idx:])
		break
	case config.APPEND:
		buffer.WriteString(fileContent)
		buffer.WriteString(markdownResult)
		break
	default:
		buffer.WriteString(fileContent)
		buffer.WriteString(markdownResult)
		break
	}

	return buffer.String(), nil
}

func (m *MarkdownGenerator) SetChangeLogHandler(clh change.ChangeLogBuilder) {
	m.changeLogHandler = clh
}

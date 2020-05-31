package change

import (
	"changelog-generator/config"
	"changelog-generator/handler/version"
	"changelog-generator/scm"
	"encoding/json"
	"path/filepath"
)

type GitlabChangeLogHandler struct {
	gitlabScmObject scm.ScmAction
	logDirectory    string
	changeLogs      map[string][]DefaultChangeLog
	versionHandler  version.VersionNumber
}

func (g *GitlabChangeLogHandler) init(c *config.Config) error {
	var err error
	g.gitlabScmObject, err = scm.GetScmHandler(c)
	g.logDirectory = c.GetChangeLogPath()
	g.versionHandler = version.VersionHandlerMap[c.GetVersionAcquisitionPolicy()]
	g.changeLogs = make(map[string][]DefaultChangeLog)
	g.versionHandler.Init(c)

	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabChangeLogHandler) collectLogs() error {
	scmFiles, err := g.gitlabScmObject.GetFiles(g.logDirectory)
	if err != nil {
		return err
	}

	for _, scmFile := range scmFiles {
		changeLog := &DefaultChangeLog{
			fileName:  filepath.Base(scmFile.FileName),
			extension: filepath.Ext(scmFile.FilePath),
			path:      scmFile.FilePath,
		}

		if scmFile.FileContent == "" {
			continue
		}

		if parseErr := json.Unmarshal([]byte(scmFile.FileContent), changeLog); parseErr != nil {
			return parseErr
		}

		if changeLog.Version != "" {
			if _, ok := g.changeLogs[changeLog.Version]; ok {
				g.changeLogs[changeLog.Version] = append(g.changeLogs[changeLog.Version], *changeLog)
			} else {
				g.changeLogs[changeLog.Version] = []DefaultChangeLog{*changeLog}
			}
			return nil
		}

		if version, versionErr := g.versionHandler.GetVersion(changeLog.GetFileName()); versionErr != nil {
			return versionErr
		} else {
			if _, ok := g.changeLogs[version]; ok {
				g.changeLogs[version] = append(g.changeLogs[version], *changeLog)
			} else {
				g.changeLogs[version] = []DefaultChangeLog{*changeLog}
			}
		}
	}

	return nil
}

func (g *GitlabChangeLogHandler) GetChangeLogInfo() map[string][]DefaultChangeLog {
	return g.changeLogs
}

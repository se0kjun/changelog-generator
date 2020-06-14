package change

import (
	"changelog-generator/config"
	changelog_err "changelog-generator/errors"
	"changelog-generator/handler/version"
	"changelog-generator/scm"
	"encoding/json"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type GitlabChangeLogHandler struct {
	gitlabScmObject scm.ScmAction
	logDirectory    string
	changeLogs      map[string][]DefaultChangeLog
	versionHandler  version.VersionNumber
}

func (g *GitlabChangeLogHandler) init(c *config.Config) error {
	var err error
	if _, ok := version.VersionHandlerMap[c.GetVersionAcquisitionPolicy()]; !ok {
		return changelog_err.NOT_FOUND_VERSION_TYPE
	}

	if g.gitlabScmObject, err = scm.GetScmHandler(c); err != nil {
		return err
	}
	g.logDirectory = c.GetChangeLogPath()
	g.versionHandler = version.VersionHandlerMap[c.GetVersionAcquisitionPolicy()]
	g.changeLogs = make(map[string][]DefaultChangeLog)
	if err = g.versionHandler.Init(c); err != nil {
		return err
	}

	log.Infof("changelog path: %s", g.logDirectory)
	return nil
}

func (g *GitlabChangeLogHandler) collectLogs() error {
	log.Infof("gitlab changelog handler is collecting logs")
	scmFiles, err := g.gitlabScmObject.GetFiles(g.logDirectory)
	if err != nil {
		log.Errorf("Cannot find following file: %s", g.logDirectory)
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
			log.Errorf("json unmarshal error: %s", parseErr.Error())
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
			log.Errorf("Getting version number error: %s", changeLog.GetFileName())
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

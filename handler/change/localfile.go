package change

import (
	"changelog-generator/config"
	"changelog-generator/handler/version"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type LocalFileChangeLogHandler struct {
	ChangeLogs     map[string][]DefaultChangeLog
	versionHandler version.VersionNumber
	logDirectory   string
}

func (c *LocalFileChangeLogHandler) init(conf *config.Config) error {
	c.logDirectory = conf.GetChangeLogPath()
	c.versionHandler = version.VersionHandlerMap[conf.GetVersionAcquisitionPolicy()]
	c.ChangeLogs = make(map[string][]DefaultChangeLog)
	c.versionHandler.Init(conf)

	log.Infof("changelog path: %s", c.logDirectory)
	return nil
}

func (c *LocalFileChangeLogHandler) collectLogs() error {
	log.Infof("local file handler is collecting logs")
	err := filepath.Walk(c.logDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("file walk error: %s", err.Error())
			return err
		}

		if info.IsDir() {
			return nil
		}

		if data, fileErr := ioutil.ReadFile(path); fileErr != nil {
			return fileErr
		} else {
			changeLog := &DefaultChangeLog{
				fileName:  filepath.Base(path),
				extension: filepath.Ext(path),
				path:      path,
			}

			if parseErr := json.Unmarshal(data, changeLog); parseErr != nil {
				return parseErr
			}

			if changeLog.Version != "" {
				if _, ok := c.ChangeLogs[changeLog.Version]; ok {
					c.ChangeLogs[changeLog.Version] = append(c.ChangeLogs[changeLog.Version], *changeLog)
				} else {
					c.ChangeLogs[changeLog.Version] = []DefaultChangeLog{*changeLog}
				}
				return nil
			}

			if version, versionErr := c.versionHandler.GetVersion(changeLog.GetFileName()); versionErr != nil {
				return versionErr
			} else {
				if _, ok := c.ChangeLogs[version]; ok {
					c.ChangeLogs[version] = append(c.ChangeLogs[version], *changeLog)
				} else {
					c.ChangeLogs[version] = []DefaultChangeLog{*changeLog}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *LocalFileChangeLogHandler) GetChangeLogInfo() map[string][]DefaultChangeLog {
	return c.ChangeLogs
}

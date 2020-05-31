package change

import (
	"changelog-generator/config"
	"changelog-generator/handler/version"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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

	return nil
}

func (c *LocalFileChangeLogHandler) collectLogs() error {
	err := filepath.Walk(c.logDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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

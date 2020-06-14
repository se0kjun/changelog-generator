package version

import (
	"changelog-generator/config"
	changelog_err "changelog-generator/errors"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type VersionNumber interface {
	Init(c *config.Config) error
	GetVersion(interface{}) (string, error)
}

var VersionHandlerMap = map[string]VersionNumber{
	// file name
	config.VERSION_GET_FILENAME: &FileBasedVersionNumber{},
	// git tag
	config.VERSION_GET_SCM_TAG: &TagBasedVersionNumber{},
}

/* FileBasedVersionNumber */
type FileBasedVersionNumber struct {
	filenameRule         string
	filenameRegexCompile *regexp.Regexp
}

func (f *FileBasedVersionNumber) Init(c *config.Config) error {
	log.Infof("filename based version initializing")
	var err error
	f.filenameRule = c.GetVersionParsingRule()
	f.filenameRegexCompile, err = regexp.Compile(f.filenameRule)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBasedVersionNumber) GetVersion(item interface{}) (string, error) {
	switch tmp := item.(type) {
	case string:
		if res := f.filenameRegexCompile.FindString(tmp); res == "" {
			return "", changelog_err.FAIL_TO_GET_VERSION
		} else {
			return res, nil
		}
	default:
		return "", changelog_err.FAIL_TO_GET_VERSION
	}
}

/* TagBasedVersionNumber */
type TagBasedVersionNumber struct {
}

func (f *TagBasedVersionNumber) Init(c *config.Config) error {
	return nil
}

func (f *TagBasedVersionNumber) GetVersion(item interface{}) (string, error) {
	return "", nil
}

package version

import (
	"changelog-generator/config"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type VersionNumber interface {
	Init(c *config.Config) error
	GetVersion(interface{}) (string, error)
}

var VersionHandlerMap = map[string]VersionNumber{
	// command line
	"default": &DefaultVersionNumber{},
	// file name
	"filename": &FileBasedVersionNumber{},
	// git tag
	"scm-tag": &TagBasedVersionNumber{},
}

/* DefaultVersionNumber */
type DefaultVersionNumber struct {
}

func (f *DefaultVersionNumber) Init(c *config.Config) error {
	return nil
}

func (f *DefaultVersionNumber) GetVersion(item interface{}) (string, error) {
	return "", nil
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
	changelog := item.(string)

	return f.filenameRegexCompile.FindString(changelog), nil
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

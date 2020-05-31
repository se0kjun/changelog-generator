package markdown

import (
	"changelog-generator/config"
	"changelog-generator/handler/change"
	"io/ioutil"
	"testing"
)

func TestGenerateMarkdown(t *testing.T) {
	c := &config.Config{
		ProjectAccessType: "localfile",
		ChangeLogConfig: config.ChangeLogGenerateConfig{
			ChangeLogPath:            "/Users/seokjunhong/Documents/golang/changelog-generator/test/unreleased",
			VersionParseRule:         "^\\d+\\.\\d+\\.\\d",
			VersionAcquisitionPolicy: "filename",
			ChangeLogOnly:            true,
		},
		OutputConfig: config.OutputConfig{
			OutputFilePath: "/Users/seokjunhong/Documents/golang/changelog-generator/test/CHANGELOG.md",
			WritePolicy:    config.APPEND,
		},
	}
	if handler, err := change.NewChangeLogHandler(c); err != nil {
		t.Error(err)
	} else {
		if markdownGen, err := NewMarkdownGenerator(c, handler); err != nil {
			t.Error(err)
		} else {
			if str, err := markdownGen.MakeResult(); err != nil {
				t.Error(err)
			} else {
				ioutil.WriteFile("/Users/seokjunhong/Documents/golang/changelog-generator/test/CHANGELOG.md", []byte(str), 0644)
				t.Log(str)
			}
		}
	}
}

package change

import (
	"changelog-generator/config"
	"testing"
)

func TestChangeLog(t *testing.T) {
	c := &config.Config{
		ProjectAccessType: "localfile",
		ChangeLogConfig: config.ChangeLogGenerateConfig{
			ChangeLogPath:            "/Users/seokjunhong/Documents/golang/changelog-generator/test/unreleased",
			VersionParseRule:         "^\\d+\\.\\d+\\.\\d",
			VersionAcquisitionPolicy: "filename",
		},
	}
	if handler, err := NewChangeLogHandler(c); err != nil {
		t.Error(err)
	} else {
		t.Log(handler.GetChangeLogInfo())
	}
}

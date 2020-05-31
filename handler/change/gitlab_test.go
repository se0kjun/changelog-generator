package change

import (
	"changelog-generator/config"
	"testing"
)

func TestGitlabMakeLog(t *testing.T) {
	c := &config.Config{
		ProjectAccessType: "gitlab",
		ChangeLogConfig: config.ChangeLogGenerateConfig{
			ChangeLogPath:            "test/unreleased",
			VersionParseRule:         "^\\d+\\.\\d+\\.\\d",
			VersionAcquisitionPolicy: "filename",
		},
		ScmConfig: config.ScmConfig{
			ScmType:           "gitlab",
			ScmAccessToken:    "<secret_access_token>",
			ScmApiBaseUrl:     "https://gitlab.com/api/v4/",
			ScmRepositoryInfo: "project_with_namespace",
			ScmPostAction: config.PostActionConfig{
				TargetBranch:  "master",
				AuthorEmail:   "hong921122@gmail.com",
				AuthorName:    "Seokjun Hong",
				CommitMessage: "Create README",
			},
		},
	}

	logBuilder, err := NewChangeLogHandler(c)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(logBuilder.GetChangeLogInfo())
	}
}

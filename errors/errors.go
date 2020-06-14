package errors

import (
	"errors"
)

var (
	NOT_FOUND_ACCESS_TYPE  = errors.New("Can't find project access type")
	NOT_FOUND_VERSION_TYPE = errors.New("Can't find version type")
	NOT_FOUND_PATH         = errors.New("Can't find following path: %s")
	SCM_NOT_INITIALIZED    = errors.New("SCM client has not initialized yet")
	FAIL_TO_GET_VERSION    = errors.New("fail to get version number")
	UNKNOWN_PROJECT_TYPE   = errors.New("Unknown project type")
	NOT_FOUND_DELIMETER    = errors.New("Not found delimeter '---'")
)

package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"strings"
	"service/service/info"
)

var Env string
var BaseAdministratorPortalUrl = fmt.Sprintf("https://admin.%s.co.za", AppProject)
var BaseApiUrl = fmt.Sprintf("https://api.%s.co.za", AppProject)
var FromEmailAddress = fmt.Sprintf("noreply@%s.co.za", AppProject)
var FromEmailName = strings.ToTitle(AppProject)

func RunBefore(p diary.IPage) {
	if value, exist := args["env"]; exist && value != nil && value != "" {
		Env = fmt.Sprint(value)
	} else {
		panic("env was not set")
	}

	// setup configurations based on running environment
	switch strings.ToLower(Env) {
	default:
		panic(fmt.Sprintf("unknown environment `%s` given", Env))
	case "prod":
		// defaults are already set to run on local environment
		break
	case "demo":
		break
	case "staging":
		break
	case "qa":
		break
	case "dev":
		break
	case "test":
		// used for automated testing sandbox environment
		break
	case "local":
		break
	}
}

var EnvPrefix = func() string {
	switch strings.ToLower(Env) {
	case "demo":
		return "[DEMO] "
	case "staging":
		return "[STAGING] "
	case "qa":
		return "[QA] "
	case "dev":
		return "[DEV] "
	case "test":
		return "[TEST] "
	case "local":
		return "[LOCAL] "
	}
	return ""
}

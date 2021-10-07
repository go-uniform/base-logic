package info

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"strings"
)

/* Details
Your hardcoded global service instance details for runtime
*/

import (
	"time"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
	AppService = "service"
	Database = AppProject
	JwtExpiryTime = time.Hour * 24
)

var Args uniform.M
var Env string
var BaseAdministratorPortalUrl = fmt.Sprintf("https://admin.%s.co.za", AppProject)
var BaseApiUrl = fmt.Sprintf("https://api.%s.co.za", AppProject)
var FromEmailAddress = fmt.Sprintf("noreply@%s.co.za", AppProject)
var FromEmailName = strings.ToTitle(AppProject)
var Salt = `9xwY7QhIfmekHQ*@qGYT#!SU9ngVdFAU`

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
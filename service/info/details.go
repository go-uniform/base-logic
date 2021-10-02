package info

import "github.com/go-uniform/uniform"

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
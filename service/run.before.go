package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"strings"
	"service/service/info"
	"sync"
)

func RunBefore(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	if value, exist := info.Args["env"]; exist && value != nil && value != "" {
		info.Env = fmt.Sprint(value)
	} else {
		panic("env was not set")
	}

	// setup configurations based on running environment
	switch strings.ToLower(info.Env) {
	default:
		panic(fmt.Sprintf("unknown environment `%s` given", info.Env))
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

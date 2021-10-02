package service

import (
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"service/service/_base"
	"service/service/info"

	// load all actions, commands, events and hooks
	_ "service/service/actions"
	_ "service/service/commands"
	_ "service/service/events"
	_ "service/service/hooks"
)

// wrap base execute to avoid circular reference
func Execute(level string, rate, limit int, test bool, natsUri string, natsOptions []nats.Option, argsMap uniform.M) {
	info.Args = argsMap
	if info.Args == nil {
		info.Args = uniform.M{}
	}
	_base.InitializeDiary(test, level, rate)
	_base.Execute(limit, test, natsUri, natsOptions, RunBefore, RunAfter)
}
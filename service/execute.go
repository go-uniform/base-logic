package service

import (
	"github.com/nats-io/go-nats"
	"service/service/_base"
)

// wrap args into service layer
var args M

// wrap base execute to avoid circular reference
func Execute(level string, rate, limit int, test bool, natsUri string, natsOptions []nats.Option, argsMap M) {
	args = argsMap
	_base.InitializeDiary(test, level, rate)
	_base.Execute(limit, test, natsUri, natsOptions, _base.M(argsMap), RunBefore, RunAfter)
}
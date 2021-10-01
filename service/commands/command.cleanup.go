package commands

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetCommand("cleanup"), cleanup)
}

func cleanup(r uniform.IRequest, p diary.IPage) {
	// todo: add routines to cleanup/migrate legacy data issues
}

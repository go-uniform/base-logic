package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(command("cleanup"), cleanup)
}

func cleanup(r uniform.IRequest, p diary.IPage) {
	// todo: add routines to cleanup/migrate legacy data issues
}

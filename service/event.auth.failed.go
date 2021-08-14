package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "failed"), eventAuthFailed)
}

func eventAuthFailed(r uniform.IRequest, p diary.IPage) {
	// todo: handle account blocking rules here
}
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "check"), eventAuthCheck)
}

func eventAuthCheck(r uniform.IRequest, p diary.IPage) {
	// todo: find account record and check rules like active/blocked, return record if auth may be attempted against the account
}
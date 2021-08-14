package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "otp"), eventAuthJwt)
}

func eventAuthJwt(r uniform.IRequest, p diary.IPage) {
	// todo: compile and return the Json-Web-Token content here
}
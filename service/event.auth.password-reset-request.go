package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "password-reset-request"), eventAuthPasswordResetRequest)
}

func eventAuthPasswordResetRequest(r uniform.IRequest, p diary.IPage) {
	// todo: create a password reset token for the given account
}
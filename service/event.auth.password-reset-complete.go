package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "password-reset-complete"), eventAuthPasswordResetComplete)
}

func eventAuthPasswordResetComplete(r uniform.IRequest, p diary.IPage) {
	// todo: complete the reset request by setting the new password for the given account and wiping the reset token
}
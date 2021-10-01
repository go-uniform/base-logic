package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "password-reset-complete"), eventAuthPasswordResetComplete)
}

func eventAuthPasswordResetComplete(r uniform.IRequest, p diary.IPage) {
	// todo: complete the reset request by setting the new password for the given account and wiping the reset token
}
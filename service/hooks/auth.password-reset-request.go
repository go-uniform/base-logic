package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "password-reset-request"), eventAuthPasswordResetRequest)
}

func eventAuthPasswordResetRequest(r uniform.IRequest, p diary.IPage) {
	// todo: create a password reset token for the given account
}
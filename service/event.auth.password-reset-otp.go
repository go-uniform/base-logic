package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "password-reset-otp"), eventAuthPasswordResetOtp)
}

func eventAuthPasswordResetOtp(r uniform.IRequest, p diary.IPage) {
	// todo: send the One-Time-Pin to the customer to complete two-factor authentication or send email containing reset link
}
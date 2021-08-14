package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("auth", "otp"), eventAuthOtp)
}

func eventAuthOtp(r uniform.IRequest, p diary.IPage) {
	// todo: send the One-Time-Pin to the customer to complete two-factor authentication
}
package actions

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

type AuthLoginRequest struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthLoginResponse struct {
	TwoFactor bool `json:"two-factor"`
	Token string `json:"token"`
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.login"), func(r uniform.IRequest, p diary.IPage) {
		var request AuthLoginRequest
		r.Read(&request)

		p.Notice("auth.login", diary.M{
			"type":       request.Type,
			"identifier": request.Identifier,
		})
		// todo: check service health and throw error if unhealthy

		if r.CanReply() {
			if err := r.Reply(uniform.Request{}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}
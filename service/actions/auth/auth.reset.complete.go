package auth

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

type ResetCompleteRequest struct {
	Type     string `bson:"type"`
	Token    string `bson:"token"`
	Password string `bson:"password"`
}

type ResetCompleteResponse struct {
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.reset.complete"), func(r uniform.IRequest, p diary.IPage) {
		var request ResetCompleteRequest
		r.Read(&request)

		p.Notice("auth.reset.complete", diary.M{
			"type":       request.Type,
			"resetToken": request.Token,
			"password":   request.Password,
		})

		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "reset.send")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"type":       request.Type,
				"resetToken": request.Token,
				"password":   request.Password,
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
		}); err != nil {
			p.Notice("failed.send", diary.M{
				"type":       request.Type,
				"resetToken": request.Token,
				"password":   request.Password,
			})
			panic(err)
		}

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: ResetCompleteResponse{
				},
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}

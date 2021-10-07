package auth

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

type LoginResendRequest struct {
	RequestId *string `bson:"token"`
}

type LoginResendResponse struct {
	RequestId *string `bson:"token"`
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.login.resend"), func(r uniform.IRequest, p diary.IPage) {
		var request LoginResendRequest
		r.Read(&request)

		p.Notice("auth.resend", diary.M{
			"requestId": request.RequestId,
		})

		var codeEntity struct {
			Id   string                 `bson:"id"`
			Code string                 `bson:"code"`
			Meta map[string]interface{} `bson:"meta"`
		}
		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("code", "reissue")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"id": request.RequestId,
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
			r.Read(&codeEntity)
		}); err != nil {
			p.Notice("failed.code", diary.M{
				"requestId": request.RequestId,
			})
			panic(err)
		}

		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "login.send")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"type": codeEntity.Meta["type"],
				"identifier": codeEntity.Meta["identifier"],
				"method": codeEntity.Meta["method"],
				"code": codeEntity.Code,
				"channel": codeEntity.Meta["channel"],
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
			r.Read(&codeEntity)
		}); err != nil {
			p.Notice("failed.send", diary.M{
				"requestId": request.RequestId,
			})
			panic(err)
		}

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: LoginResendResponse{
					RequestId: &codeEntity.Id,
				},
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}

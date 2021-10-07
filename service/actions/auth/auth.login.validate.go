package auth

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

type LoginValidateRequest struct {
	Code  string `bson:"code"`
	RequestId *string `bson:"token"`
}

type LoginValidateResponse struct {
	Token string `bson:"token"`
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.login.validate"), func(r uniform.IRequest, p diary.IPage) {
		var request LoginValidateRequest
		r.Read(&request)

		p.Notice("auth.login.validate", diary.M{
			"code": request.Code,
			"requestId": request.RequestId,
		})

		var codeEntity struct{
			Id string `bson:"id"`
			Code string `bson:"code"`
			Meta map[string]interface{} `bson:"meta"`
		}
		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("code", "consume")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"id": request.RequestId,
				"code": request.Code,
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

		token := ""
		if val, exists := codeEntity.Meta["token"]; exists {
			if strVal, ok := val.(string); ok {
				token = strVal
			}
		}

		if token == "" {
			panic("unable to retrieve account login token")
		}

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: LoginValidateResponse{
					Token: token,
				},
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}

package auth

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"service/service/_base"
	"strings"
	"time"
)

type ResetRequest struct {
	Group      string `json:"group"`
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
	Method     string `json:"method"`
	Channel    string `json:"channel"`

	Id   string                 `json:"id"`
	Meta map[string]interface{} `json:"meta"`
}

type ResetResponse struct {
	RequestId *string `json:"request-id"`
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.reset"), func(r uniform.IRequest, p diary.IPage) {
		var request ResetRequest
		r.Read(&request)

		p.Notice("auth.reset", diary.M{
			"type":       request.Type,
			"identifier": request.Identifier,
		})
		id := request.Id
		if id == "" {
			var entity Entity
			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "check")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(&entity)
			}); err != nil {
				p.Notice("failed.auth.check", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}

			id = entity.Id
		}

		var token = primitive.NewObjectID()
		var requestId *string

		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "reset.request")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"meta": uniform.M{
					"type":                request.Type,
					"id":                  id,
					"resetToken":          token,
					"resetTokenExpiresAt": time.Now().Add(time.Hour * 3),
				},
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
		}); err != nil {
			p.Notice("failed.reset.request", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		if strings.ToLower(request.Method) == "code" {
			var codeEntity struct {
				Id   string `json:"id"`
				Code string `json:"code"`
			}
			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("code", "issue")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"meta": uniform.M{
						"token":      token,
						"type":       request.Type,
						"identifier": request.Identifier,
						"method":     request.Method,
						"channel":    request.Channel,
						"id":         request.Id,
						"meta":       request.Meta,
					},
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(&codeEntity)
			}); err != nil {
				p.Notice("failed.code", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}

			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "reset.send")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"type":       request.Type,
					"identifier": request.Identifier,
					"method":     "code",
					"code":       codeEntity.Code,
					"channel":    request.Channel,
					"id":         request.Id,
					"meta":       request.Meta,
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(&codeEntity)
			}); err != nil {
				p.Notice("failed.send", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}

			requestId = &codeEntity.Id
		} else if strings.ToLower(request.Method) == "token" {
			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "reset.send")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"type":       request.Type,
					"identifier": request.Identifier,
					"method":     "token",
					"token":      token.Hex(),
					"channel":    request.Channel,
					"id":         request.Id,
					"meta":       request.Meta,
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
			}); err != nil {
				p.Notice("failed.send", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}
		} else {
			panic("unsupported method")
		}

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: ResetResponse{
					RequestId: requestId,
				},
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}

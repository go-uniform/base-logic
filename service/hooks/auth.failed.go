package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
	"strings"
	"time"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "failed"), eventAuthFailed)
}

func eventAuthFailed(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthFailedRequest
	var response uniform.AuthFailedResponse
	r.Read(&request)

	db := r.Conn().Mongo(p, "")
	exists := false
	db.CatchNoDocumentsErr(func(p diary.IPage) {
		switch strings.ToLower(request.Type) {
		default:
			p.Warning("check", "an attempt to auth an unknown type", diary.M{
				"id":   request.Id,
				"type": request.Type,
			})
			uniform.Alert(401, "Incorrect login details")
		case "administrator":
			db.Read(r.Remainder(), _base.Database, "administrators", request.Id, &response, nil)
			if response.LockedAt == nil && response.BlockedAt == nil {
				if response.Counter >= 2 {
					db.Update(r.Remainder(), _base.Database, "administrators", request.Id, uniform.M{
						"lockedAt": time.Now(),
					}, nil, nil)
				} else {
					db.Inc(r.Remainder(), _base.Database, "administrators", request.Id, "counter", 1, nil, nil)
				}
			}
			break
		}
		exists = true
	})

	if err := r.Reply(uniform.Request{
		Model: response,
	}); err != nil {
		p.Error("reply", err.Error(), diary.M{
			"error": err,
		})
	}
}

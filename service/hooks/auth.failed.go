package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/go-uniform/uniform/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"service/service/_base"
	"service/service/info"
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

	db := nosql.Request(r.Conn(), p, "")
	exists := false
	db.CatchErrNoResults(func(p diary.IPage) {
		switch strings.ToLower(request.Type) {
		default:
			p.Warning("check", "an attempt to auth an unknown type", diary.M{
				"id":   request.Id,
				"type": request.Type,
			})
			uniform.Alert(401, "Incorrect login details")
		case "administrator":
			db.FindOne(r.Remainder(), info.Database, "administrators", "", 0, bson.D{
				{"_id", request.Id},
			}, &response)
			if response.LockedAt == nil && response.BlockedAt == nil {
				if response.Counter >= 2 {
					now := time.Now()
					response.LockedAt = &now
				} else {
					response.Counter++
				}
				db.UpdateOne(r.Remainder(), info.Database, "administrators", request.Id, response, nil)
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

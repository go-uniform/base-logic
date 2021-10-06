package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/go-uniform/uniform/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"service/service/_base"
	"service/service/info"
	"strings"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "check"), eventAuthCheck)
}

func eventAuthCheck(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthCheckRequest
	var response uniform.AuthCheckResponse
	r.Read(&request)

	db := nosql.Request(r.Conn(), p, "")
	exists := false
	db.CatchErrNoResults(func(p diary.IPage) {
		switch strings.ToLower(request.Type) {
		default:
			p.Warning("check", "an attempt to auth an unknown type", diary.M{
				"identifier": request.Identifier,
				"type":       request.Type,
			})
			uniform.Alert(401, "Incorrect login details")
		case "administrator":
			db.FindOne(r.Remainder(), info.Database, "administrators", "", 0, bson.D{
				{"_id", request.Identifier },
			}, &response)
			break
		}
		exists = true
	})

	if !exists {
		// keep in mind that because of this error people may use the login function to find out if another person is a member of your site
		// so if your site has an element of confidentiality associated with it you will want to remove this error to keep member's existence confidential
		// having this check does help user experience specifically for typos in identifier so don't remove if you don't have to
		uniform.Alert(401, "On account matched the given identifier")
	}

	if response.BlockedAt != nil {
		uniform.Alert(403, "Account has been blocked by an administrator")
	}

	if !request.Reset {
		if response.Password == nil {
			uniform.Alert(403, "Account not yet activated, please do a password reset")
		} else if response.LockedAt != nil {
			uniform.Alert(403, "Account has been locked because of too many failed login attempts, please do a password reset")
		}
	}

	if err := r.Reply(uniform.Request{
		Model: response,
	}); err != nil {
		p.Error("reply", err.Error(), diary.M{
			"error": err,
		})
	}
}

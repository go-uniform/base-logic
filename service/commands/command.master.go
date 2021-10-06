package commands

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/go-uniform/uniform/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
	"time"
)

func init() {
	_base.Subscribe(_base.TargetCommand("master"), func(r uniform.IRequest, p diary.IPage) {
		var request struct {
			FirstName string `json:"firstName"`
			LastName string `json:"lastName"`
			Email string `json:"email"`
			Mobile string `json:"mobile"`
			Password string `json:"password"`
		}
		r.Read(&request)

		db := nosql.Request(r.Conn(), p, "")
		if db.Count(time.Second * 5, info.Database, entities.CollectionAdministrators, bson.D{}) > 0 {
			if r.CanReply() {
				if err := r.Reply(uniform.Request{
					Error: "The collection already contains records so a master record can't be created at this point",
				}); err != nil {
					p.Error("reply", err.Error(), diary.M{
						"error": err,
					})
				}
			}
			return
		}

		db.InsertOne(time.Second * 5, info.Database, "administrators", entities.Administrator{
			FirstName: request.FirstName,
			LastName: request.LastName,
			Email: request.Email,
			Mobile: request.Mobile,
			Password: request.Password,
			Inverted: true,
		}, nil)

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: "Created master record!",
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"error": err,
				})
			}
		}
	})
}

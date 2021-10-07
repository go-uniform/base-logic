package entity

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"go.mongodb.org/mongo-driver/bson"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
)

func init() {
	_base.Subscribe(_base.TargetEvent("entity", fmt.Sprintf("%s.validate", info.Database)), func(r uniform.IRequest, p diary.IPage) {
		var response interface{}

		switch r.Parameters()["collection"] {
		default:
			// send request straight to response
			r.Read(&response)
			break

		case entities.CollectionAdministrators:
			var entity bson.M
			r.Read(&entity)

			// set the identifier
			entity["identifier"] = []string{
				uniform.Hash(entity["email"], info.Salt),
			}

			response = entity
			break

		}

		if err := r.Reply(uniform.Request{
			Model: response,
		}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"error": err,
			})
		}
	})
}

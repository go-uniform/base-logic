package entity

import (
	"encoding/base64"
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
)

func init() {
	_base.Subscribe(_base.TargetEvent("entity", fmt.Sprintf("%s.encrypt", info.Database)), func(r uniform.IRequest, p diary.IPage) {
		var response interface{}

		switch r.Parameters()["collection"] {
		default:
			// send request straight to response
			r.Read(&response)

		case entities.CollectionAdministrators:
			var entity bson.M
			r.Read(&entity)

			hashed, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprint(entity["password"])), bcrypt.MinCost)
			if err != nil {
				panic(err)
			}
			entity["password"] = base64.StdEncoding.EncodeToString(hashed)

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

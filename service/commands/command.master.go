package commands

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service"
	"service/service/_base"
	"time"
)

func init() {
	_base.Subscribe(_base.TargetCommand("master"), master)
}

func master(r uniform.IRequest, p diary.IPage) {
	var request struct {
		FirstName string `json:"first-name"`
		LastName string `json:"last-name"`
		Email string `json:"email"`
		Mobile string `json:"mobile"`
	}
	r.Read(&request)

	db := r.Conn().Mongo(p, "")
	if db.Count(time.Second * 5, service.AppService, "administrators", uniform.M{}) > 0 {
		panic("The collection already contains records so a master record can't be created at this point")
	}

	db.Insert(time.Second * 5, service.AppService, "administrators", service.Administrator{
		FirstName: request.FirstName,
		LastName: request.LastName,
		Email: request.Email,
		Mobile: request.Mobile,
	}, nil, service.TagsAdministrator)

	if r.CanReply() {
		if err := r.Reply(uniform.Request{
			Model: "Created master record!",
		}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"error": err,
			})
		}
	}
}

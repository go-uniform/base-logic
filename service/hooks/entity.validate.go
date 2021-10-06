package hooks

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
	"service/service/info"
)

func init() {
	_base.Subscribe(_base.TargetEvent("entity", fmt.Sprintf("%s.validate", info.Database)), func(r uniform.IRequest, p diary.IPage) {
		var response interface{}

		switch r.Parameters()["collection"] {
		default:
			// send request straight to response
			r.Read(&response)

		// todo: add validation rules

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

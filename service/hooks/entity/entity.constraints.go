package entity

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
	"service/service/info"
)

func init() {
	_base.Subscribe(_base.TargetEvent("entity", fmt.Sprintf("%s.constraints", info.Database)), func(r uniform.IRequest, p diary.IPage) {
		var response interface{}

		switch r.Parameters()["collection"] {

			// todo: add constaint checks

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

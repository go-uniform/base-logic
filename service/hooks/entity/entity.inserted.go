package entity

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetEvent("entity", "inserted"), entityInserted)
}

func entityInserted(r uniform.IRequest, p diary.IPage) {
	// todo: react based on which database and collection has been inserted
}
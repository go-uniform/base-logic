package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local(command("example-one")), exampleOne)
}

func exampleOne(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}

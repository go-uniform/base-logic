package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-three"), exampleThree)
}

func exampleThree(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
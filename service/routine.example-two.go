package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-two"), exampleTwo)
}

func exampleTwo(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("mongo", "inserted"), eventMongoInserted)
}

func eventMongoInserted(r uniform.IRequest, p diary.IPage) {
	// todo: react based on which database and collection has been inserted
}
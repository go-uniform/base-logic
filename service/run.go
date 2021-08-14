package service

import (
	"github.com/go-diary/diary"
	"time"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
	AppService = "service"
	Database = AppProject
	DatabaseTimeout = time.Second * 5
)

func Run(p diary.IPage) {
	// todo: add your custom start up logic here
}
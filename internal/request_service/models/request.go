package models

import (
	"github.com/google/uuid"
)

type Request struct {
	Id         uuid.UUID
	URLBuilder string
	BuildArgs  []interface{}
	JITBuild   func(urlBuilder string, buildArgs []interface{}) (string, error)
}

type RequestResponse struct {
	Id   uuid.UUID
	Resp []byte
	Err  error
}

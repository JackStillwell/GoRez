package models

import (
	"github.com/google/uuid"
)

type Request struct {
	Id      *uuid.UUID
	JITFunc func() (string, error)
}

type RequestResponse struct {
	Id   *uuid.UUID
	Resp []byte
	Err  error
}

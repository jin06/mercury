package dao

import "github.com/jin06/mercury/internal/model"

type Record interface{}

type memRecord struct {
	used []*model.Record
}

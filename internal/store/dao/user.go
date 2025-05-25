package dao

import "github.com/jin06/mercury/internal/model"

type User struct {
}

func (d *User) Get(name string) (*model.User, error) {
	return nil, nil
}

func (d *User) Create(account *model.User) error {
	return nil
}

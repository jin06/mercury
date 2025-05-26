package dao

import (
	"github.com/jin06/mercury/internal/model"
	"gorm.io/gorm"
)

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

type User struct {
	db *gorm.DB
}

func (d *User) Get(name string) (*model.User, error) {
	return nil, nil
}

func (d *User) Create(data *model.User) error {
	return d.db.Create(data).Error
}

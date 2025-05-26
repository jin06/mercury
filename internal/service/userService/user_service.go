package userservice

import (
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/store"
	"github.com/jin06/mercury/internal/store/dao"
)

var Single *UserService

func Init() {
	Single = NewUserService()
}

func Get(name string) (*model.User, error) {
	return Single.Get(name)
}

func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUser(store.Default),
	}
}

type UserService struct {
	userDao *dao.User
}

func (s *UserService) Get(name string) (*model.User, error) {
	user, err := s.userDao.Get(name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

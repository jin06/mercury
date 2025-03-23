package dao

type Account interface{}

type dbAccount struct {
}

func (d *dbAccount) GetAccount(name string) (*Account, error) {
	return nil, nil
}

package dao

type Account interface{}

type memAccount struct {
}

func (d *memAccount) GetAccount(name string) (*Account, error) {
	return nil, nil
}

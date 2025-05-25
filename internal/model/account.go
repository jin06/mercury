package model

import "time"

// Account represents the account information for an MQTT user.
type Account struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`           // Unique identifier for the account
	Username   string    `json:"username" gorm:"username"`       // Username
	Password   string    `json:"password" gorm:"password"`       // Password
	Created    time.Time `json:"created" gorm:"created"`         // Creation time of the account
	Updated    time.Time `json:"updated" gorm:"updated"`         // Last update time of the account
	LoginTime  time.Time `json:"login_time" gorm:"login_time"`   // Last login time of the account
	LoginCount uint64    `json:"login_count" gorm:"login_count"` // Number of times the account has logged in
}

func (a *Account) TableName() string {
	return "accounts"
}

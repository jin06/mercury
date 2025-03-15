package model

// Account represents the account information for an MQTT user.
type Account struct {
	Username string `json:"username"` // Username
	Password string `json:"password"` // Password
}

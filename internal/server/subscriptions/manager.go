package subscriptions

// todo test
type SubManager interface {
	Sub(topic string, clientID string) error
	Unsub(topic string, clientID string)
	GetSubers(topic string) []*Subscriber
}

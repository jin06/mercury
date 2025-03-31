package subscriptions

// todo test
type SubManager interface {
	Sub(topic string, clientID string) (bool, error)
	Unsub(topic string, clientID string) bool
	GetSubers(topic string) []*Subscriber
}

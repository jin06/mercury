package subscriptions

type Subscriber struct {
	Type     Type
	ClientID string
	Group    string
}

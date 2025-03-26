package subscriptions

import (
	"strings"
	"sync"
)

type trieNode struct {
	children map[string]*trieNode
	subs     map[string]*Subscriber
	mu       sync.RWMutex
}

type trieSub struct {
	root *trieNode
}

func NewTrie() *trieSub {
	return &trieSub{
		root: &trieNode{
			children: make(map[string]*trieNode),
			subs:     make(map[string]*Subscriber),
		},
	}
}

func (t *trieSub) Sub(topic string, clientID string) error {
	// parts := strings.Split(topic, "/")
	tf, err := NewTF(topic)
	if err != nil {
		return err
	}
	node := t.root
	for _, part := range tf.Parts {
		node.mu.Lock()
		if _, ok := node.children[part]; !ok {
			node.children[part] = &trieNode{
				children: make(map[string]*trieNode),
				subs:     make(map[string]*Subscriber),
			}
		}
		node.mu.Unlock()
		node = node.children[part]
	}
	node.mu.Lock()
	node.subs[clientID] = tf.subscriber(clientID)
	node.mu.Unlock()
	return nil
}

func (t *trieSub) Unsub(topic string, clientID string) {
	parts := strings.Split(topic, "/")
	node := t.root
	for _, part := range parts {
		node.mu.RLock()
		if _, ok := node.children[part]; !ok {
			node.mu.RUnlock()
			return
		}
		node = node.children[part]
		node.mu.RUnlock()
	}
	node.mu.Lock()
	delete(node.subs, clientID)
	node.mu.Unlock()
}

func (t *trieSub) GetSubers(topic string) []*Subscriber {
	parts := strings.Split(topic, "/")
	node := t.root
	for _, part := range parts {
		if found := func() bool {
			node.mu.RLock()
			defer node.mu.RUnlock()
			if _, ok := node.children[part]; !ok {
				return false
			}
			node = node.children[part]
			return true
		}(); !found {
			return nil
		}
	}
	node.mu.RLock()
	defer node.mu.RUnlock()
	subs := make([]*Subscriber, 0, len(node.subs))
	for _, suber := range node.subs {
		subs = append(subs, suber)
	}
	return subs
}

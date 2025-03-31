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

func (t *trieSub) Sub(topic string, clientID string) (bool, error) {
	tf, err := NewTF(topic)
	if err != nil {
		return false, err
	}
	node := t.root
	var has bool

	for _, part := range tf.Parts {
		node.mu.Lock()
		child, ok := node.children[part]
		if !ok {
			child = &trieNode{
				children: make(map[string]*trieNode),
				subs:     make(map[string]*Subscriber),
			}
			node.children[part] = child
		}
		node.mu.Unlock()
		node = child
	}

	node.mu.Lock()
	defer node.mu.Unlock()
	_, has = node.subs[clientID]
	node.subs[clientID] = tf.subscriber(clientID)

	return has, nil
}

func (t *trieSub) Unsub(topic string, clientID string) bool {
	parts := strings.Split(topic, "/")
	node := t.root
	var parent *trieNode
	var key string

	for _, part := range parts {
		node.mu.RLock()
		child, ok := node.children[part]
		node.mu.RUnlock()
		if !ok {
			return false
		}
		parent = node
		key = part
		node = child
	}

	node.mu.Lock()
	defer node.mu.Unlock()
	if _, exists := node.subs[clientID]; !exists {
		return false
	}
	delete(node.subs, clientID)

	// Clean up empty nodes
	if len(node.subs) == 0 && len(node.children) == 0 && parent != nil {
		parent.mu.Lock()
		delete(parent.children, key)
		parent.mu.Unlock()
	}

	return true
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

package subscriptions

import (
	"strings"
	"sync"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type RetainManager interface {
	Insert(publish *mqtt.Publish) (bool, error)
	Get(topic string) *mqtt.Publish
}

type trieNodeRetain struct {
	children map[string]*trieNodeRetain
	// subs     map[string]*model.Record
	content model.Retain
	mu      sync.RWMutex
}

type trieSubRetain struct {
	root *trieNodeRetain
}

func NewTrieRetain() *trieSubRetain {
	return &trieSubRetain{
		root: &trieNodeRetain{
			children: make(map[string]*trieNodeRetain),
			// subs:     make(map[string]*model.Record),
		},
	}
}

func (t *trieSubRetain) Insert(p *mqtt.Publish) (bool, error) {
	// tf, err := NewTF(topic)
	// return false, err
	// }
	node := t.root
	var has bool

	content := model.NewRetain(*p)
	for _, part := range p.Topic.Split() {
		node.mu.Lock()
		child, ok := node.children[part]
		if !ok {
			child = &trieNodeRetain{
				children: make(map[string]*trieNodeRetain),
				content:  model.NewRetain(*p),
			}
			node.children[part] = child
		}
		node.mu.Unlock()
		node = child
	}

	node.mu.Lock()
	defer node.mu.Unlock()
	node.content = content

	return has, nil
}

func (t *trieSubRetain) Get(topic string) *mqtt.Publish {
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
	return &node.content.Publish
}

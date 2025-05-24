package subscriptions

import (
	"strings"
	"sync"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type RetainManager interface {
	Insert(publish *mqtt.Publish) (bool, error)
	Get(topic string) []*mqtt.Publish
}

type trieNodeRetain struct {
	children map[string]*trieNodeRetain
	content  *model.Retain
	mu       sync.RWMutex
}

func (t *trieNodeRetain) GetAll() (list []*mqtt.Publish) {
	list = make([]*mqtt.Publish, 0)
	var traverse func(node *trieNodeRetain)
	traverse = func(node *trieNodeRetain) {
		node.mu.RLock()
		defer node.mu.RUnlock()
		for _, child := range node.children {
			if child.content != nil {
				list = append(list, &child.content.Publish)
			}
			traverse(child)
		}
	}
	traverse(t)
	return list
}

func (t *trieNodeRetain) Get(topic string) (list []*mqtt.Publish) {
	list = make([]*mqtt.Publish, 0)
	parts := strings.SplitN(topic, "/", 2)
	node := t
	node.mu.RLock()
	defer node.mu.RUnlock()
	if len(parts) == 0 {
		return
	}
	var p0, p1 string

	if len(parts) == 1 {
		p0 = parts[0]
		switch p0 {
		case "#":
			list = node.GetAll()
		case "+":
			for _, v := range node.children {
				if v.content != nil {
					list = append(list, &v.content.Publish)
				}
			}
		case topic:
			if child, ok := node.children[p0]; ok {
				if child.content != nil {
					list = []*mqtt.Publish{&child.content.Publish}
				}
			}
		}
		return
	}
	p1 = parts[1]
	switch p0 {
	case "#":
		panic("invalid topic: " + topic)
	case "+":
		for _, child := range node.children {
			list = append(list, child.Get(p1)...)
		}
	default:
		if _, ok := node.children[p0]; ok {
			list = node.children[p0].Get(p1)
		}
	}

	return
}

func (t *trieNodeRetain) insert(topic string, p *mqtt.Publish) (ok bool, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	parts := strings.SplitN(topic, "/", 2)
	if len(parts) == 0 {
		return false, nil
	}
	if _, ok := t.children[parts[0]]; !ok {
		t.children[parts[0]] = &trieNodeRetain{
			children: make(map[string]*trieNodeRetain),
		}
	}
	if len(parts) == 1 {
		if t.children[parts[0]].content != nil {
			ok = true
		}
		t.children[parts[0]].content = model.NewRetain(*p)
		return
	}
	return t.children[parts[0]].insert(parts[1], p)
}

func (t *trieNodeRetain) Insert(p *mqtt.Publish) (bool, error) {
	return t.insert(p.Topic.String(), p)
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

func (t *trieSubRetain) Insert(p *mqtt.Publish) (ok bool, err error) {
	return t.root.Insert(p)
}

func (t *trieSubRetain) Get(topic string) (list []*mqtt.Publish) {
	return t.root.Get(topic)
}

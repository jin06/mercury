package messages

import (
	"fmt"

	"github.com/derekparker/trie"
)

// Subscriptions pub/sub,prefix tree
// https://github.com/viant/ptrie
// https://github.com/derekparker/trie/blob/v3/trie.go
func NewTrie() {
	t := trie.New()

	node := t.Add("hello", "hello value")
	t.Add("hello/world", "hello world ")
	fmt.Println(node)
	t.Find("hello")
}

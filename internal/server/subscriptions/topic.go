package subscriptions

import (
	"errors"
	"strings"
	"time"

	"github.com/jin06/mercury/internal/utils"
)

type Type byte

const (
	TypeCommon = 0
	TypeShare  = 1
	TypeSystem = 2
)

func NewTF(rawName string) (*TopicFilter, error) {
	tf := &TopicFilter{RawName: rawName}
	if err := tf.valid(); err != nil {
		return nil, err
	}
	if err := tf.init(); err != nil {
		return nil, err
	}
	return tf, nil
}

type TopicFilter struct {
	RawName   string
	TopicName string
	Type      Type
	Group     string // share topic group
	Parts     []string
}

func (tf *TopicFilter) valid() error {
	if len(tf.RawName) == 0 {
		return utils.ErrNotValidTopic
	}
	levels := strings.Split(tf.RawName, "/")
	for i, level := range levels {
		if strings.Contains(level, "+") && len(level) != 1 {
			return utils.ErrNotValidTopic
		}
		if strings.Contains(level, "#") {
			if len(level) != 1 {
				return utils.ErrNotValidTopic
			}
			if len(levels)-1 > i {
				return utils.ErrNotValidTopic
			}
		}
	}
	return nil
}

func (tf *TopicFilter) init() error {
	tf.Parts = strings.Split(tf.RawName, "/")
	if strings.HasPrefix(tf.RawName, "$share/") {
		tf.Type = TypeShare
		parts := strings.SplitN(tf.RawName, "/", 3)
		if len(parts) < 3 {
			return errors.New("invalid shared topic format")
		}
		tf.Group = parts[1]
		tf.TopicName = parts[2]
	} else if strings.HasPrefix(tf.RawName, "$SYS/") {
		tf.Type = TypeSystem
		tf.TopicName = tf.RawName
	} else {
		tf.Type = TypeCommon
		tf.TopicName = tf.RawName
	}
	return nil
}

func (tf *TopicFilter) subscriber(clientID string) *Subscriber {
	return &Subscriber{
		Type:     tf.Type,
		ClientID: clientID,
		Group:    tf.Group,
		Time:     time.Now(),
	}
}

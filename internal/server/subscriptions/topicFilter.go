package subscriptions

import (
	"errors"
	"strings"
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
		return errors.New("topic name cannot be empty")
	}
	if strings.Contains(tf.RawName, "#") && !strings.HasSuffix(tf.RawName, "/#") {
		return errors.New("invalid use of wildcard '#' in topic name")
	}
	if strings.Contains(tf.RawName, "+") {
		segments := strings.Split(tf.RawName, "/")
		for _, segment := range segments {
			if segment == "+" && len(segment) != 1 {
				return errors.New("invalid use of wildcard '+' in topic name")
			}
		}
	}
	if strings.Contains(tf.RawName, "//") {
		return errors.New("topic name cannot contain empty levels")
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
		tf.Type,
		clientID,
		tf.Group,
	}
}
